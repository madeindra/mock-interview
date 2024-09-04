package handler

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/madeindra/mock-interview/server/internal/data"
	"github.com/madeindra/mock-interview/server/internal/middleware"
	"github.com/madeindra/mock-interview/server/internal/model"
	"github.com/madeindra/mock-interview/server/internal/openai"
	"github.com/madeindra/mock-interview/server/internal/util"
)

func (h *handler) Status(w http.ResponseWriter, _ *http.Request) {
	isKeyValid, err := h.ai.IsKeyValid()
	if err != nil {
		log.Printf("failed to check key validity: %v", err)
		util.SendResponse(w, nil, "failed to check key validity", http.StatusInternalServerError)

		return
	}

	status, err := h.ai.Status()
	if err != nil {
		log.Printf("failed to check API availability: %v", err)
		util.SendResponse(w, nil, "failed to check API availability", http.StatusInternalServerError)

		return
	}

	var apiState *bool

	switch status {
	case openai.STATUS_OPERATIONAL:
		apiState = util.Pointer(true)
	case openai.STATUS_DEGRADED_PERFORMANCE, openai.STATUS_PARTIAL_OUTAGE, openai.STATUS_MAJOR_OUTAGE:
		apiState = util.Pointer(false)
	case openai.STATUS_UNKNOWN:
		apiState = nil
	}

	apiStatus := util.Pointer(string(status))

	response := model.StatusResponse{
		Server:    true,       // always true when the server is running
		Key:       isKeyValid, // true if the API key is valid, false otherwise
		API:       apiState,   // nil if status unknown, true if operational, false otherwise
		ApiStatus: apiStatus,  // always return the status string
	}

	util.SendResponse(w, response, "success", http.StatusOK)
}

func (h *handler) StartChat(w http.ResponseWriter, req *http.Request) {
	var startChatRequest model.StartChatRequest
	if err := json.NewDecoder(req.Body).Decode(&startChatRequest); err != nil {
		log.Printf("failed to read start chat request body: %v", err)
		util.SendResponse(w, nil, "failed to read request", http.StatusBadRequest)

		return
	}

	chatLanguage := h.ai.GetDefaultTranscriptLanguage()
	if startChatRequest.Language != "" {
		chatLanguage = h.ai.GetLanguage(startChatRequest.Language)
	}

	systempPrompt, err := openai.GetSystemPrompt(startChatRequest.Role, startChatRequest.Skills, chatLanguage)
	if err != nil {
		log.Printf("failed to get system prompt: %v", err)
		util.SendResponse(w, nil, "failed to prepare chat", http.StatusInternalServerError)

		return
	}

	initialText, err := openai.GetInitialChat(startChatRequest.Role, chatLanguage)
	if err != nil {
		log.Printf("failed to get initial text: %v", err)
		util.SendResponse(w, nil, "failed to prepare chat", http.StatusInternalServerError)

		return
	}

	initialAudio, err := h.ai.TextToSpeech(util.SanitizeString(initialText))
	if err != nil {
		log.Printf("failed to create initial audio: %v", err)
		util.SendResponse(w, nil, "failed to prepare chat", http.StatusInternalServerError)

		return
	}

	audioByte, err := io.ReadAll(initialAudio)
	if err != nil {
		log.Printf("failed to read speech: %v", err)
		util.SendResponse(w, nil, "failed to read speech", http.StatusInternalServerError)

		return
	}
	audioBase64 := base64.StdEncoding.EncodeToString(audioByte)

	plainSecret := util.GenerateRandom()
	hashed, err := util.CreateHash(plainSecret)
	if err != nil {
		log.Printf("failed to create hash: %v", err)
		util.SendResponse(w, nil, "failed to prepare chat", http.StatusInternalServerError)

		return
	}

	tx, err := h.db.BeginTx()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		util.SendResponse(w, nil, "failed to create new chat", http.StatusInternalServerError)

		return
	}
	defer tx.Rollback()

	newUser, err := h.db.CreateChatUser(tx, hashed, chatLanguage)
	if err != nil {
		log.Printf("failed to create new chat: %v", err)
		util.SendResponse(w, nil, "failed to create new chat", http.StatusInternalServerError)

		return
	}

	if _, err := h.db.CreateChats(tx, newUser.ID, []data.Entry{
		{
			Role: string(openai.ROLE_SYSTEM),
			Text: systempPrompt,
		},
		{
			Role:  string(openai.ROLE_ASSISTANT),
			Text:  initialText,
			Audio: audioBase64,
		},
	}); err != nil {
		log.Printf("failed to create chat: %v", err)
		util.SendResponse(w, nil, "failed to create chat", http.StatusInternalServerError)

		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("failed to commit transaction: %v", err)
		util.SendResponse(w, nil, "failed to create new chat", http.StatusInternalServerError)

		return
	}

	initialChat := model.StartChatResponse{
		ID:       newUser.ID,
		Secret:   plainSecret,
		Language: startChatRequest.Language,
		Chat: model.Chat{
			Text:  initialText,
			Audio: audioBase64,
		},
	}

	util.SendResponse(w, initialChat, "a new chat created", http.StatusOK)
}

func (h *handler) AnswerChat(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value(middleware.ContextKeyUserID).(string)
	userSecret := req.Context().Value(middleware.ContextKeyUserSecret).(string)

	if userID == "" || userSecret == "" {
		log.Println("user ID or secret is missing")
		util.SendResponse(w, nil, "missing required authentication", http.StatusUnauthorized)

		return
	}

	user, err := h.db.GetChatUser(userID)
	if err != nil {
		log.Printf("failed to get chat user: %v", err)
		util.SendResponse(w, nil, "failed to get chat user", http.StatusNotFound)

		return
	}

	if err := util.CompareHash(userSecret, user.Secret); err != nil {
		log.Println("invalid user secret")
		util.SendResponse(w, nil, "invalid user secret", http.StatusUnauthorized)

		return
	}

	entries, err := h.db.GetChatsByChatUserID(user.ID)
	if err != nil {
		log.Printf("failed to get chat: %v", err)
		util.SendResponse(w, nil, "failed to get chat", http.StatusInternalServerError)

		return
	}

	file, fileHeader, err := req.FormFile("file")
	if err != nil {
		log.Printf("failed to read file: %v", err)
		util.SendResponse(w, nil, "failed to read file", http.StatusInternalServerError)

		return
	}
	if fileHeader == nil {
		log.Println("required file is missing")
		util.SendResponse(w, nil, "required file is missing", http.StatusBadRequest)

		return
	}
	defer file.Close()

	transcriptText, err := util.TranscribeSpeech(h.ai, file, fileHeader.Filename, user.Language)
	if err != nil {
		log.Printf("failed to transcribe speech: %v", err)
		util.SendResponse(w, nil, "failed to transcribe speech", http.StatusInternalServerError)

		return
	}

	history := util.ConvertToChatMessage(entries)

	chatHistory := append(history, openai.ChatMessage{
		Role:    openai.ROLE_USER,
		Content: transcriptText,
	})

	answerText, err := util.GenerateText(h.ai, chatHistory)
	if err != nil {
		log.Printf("failed to get chat completion: %v", err)
		util.SendResponse(w, nil, "failed to get chat completion", http.StatusInternalServerError)

		return
	}

	answerAudio, err := util.GenerateSpeech(h.ai, user.Language, answerText)
	if err != nil {
		log.Printf("failed to generate speech: %v", err)
		util.SendResponse(w, nil, "failed to generate speech", http.StatusInternalServerError)

		return
	}

	tx, err := h.db.BeginTx()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		util.SendResponse(w, nil, "failed to create new chat", http.StatusInternalServerError)

		return
	}
	defer tx.Rollback()

	if _, err := h.db.CreateChats(tx, userID, []data.Entry{
		{
			Role: string(openai.ROLE_USER),
			Text: transcriptText,
		},
		{
			Role:  string(openai.ROLE_ASSISTANT),
			Text:  answerText,
			Audio: answerAudio,
		},
	}); err != nil {
		log.Printf("failed to create chat: %v", err)
		util.SendResponse(w, nil, "failed to create chat", http.StatusInternalServerError)

		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("failed to commit transaction: %v", err)
		util.SendResponse(w, nil, "failed to create new chat", http.StatusInternalServerError)

		return
	}

	response := model.AnswerChatResponse{
		Language: h.ai.GetCode(user.Language),
		Prompt: model.Chat{
			Text: transcriptText,
		},
		Answer: model.Chat{
			Text:  answerText,
			Audio: answerAudio,
		},
	}

	util.SendResponse(w, response, "success", http.StatusOK)
}

func (h *handler) EndChat(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value(middleware.ContextKeyUserID).(string)
	userSecret := req.Context().Value(middleware.ContextKeyUserSecret).(string)

	if userID == "" || userSecret == "" {
		log.Println("user ID or secret is missing")
		util.SendResponse(w, nil, "missing required authentication", http.StatusUnauthorized)

		return
	}

	user, err := h.db.GetChatUser(userID)
	if err != nil {
		log.Printf("failed to get chat user: %v", err)
		util.SendResponse(w, nil, "failed to get chat user", http.StatusNotFound)

		return
	}

	if err := util.CompareHash(userSecret, user.Secret); err != nil {
		log.Println("invalid user secret")
		util.SendResponse(w, nil, "invalid user secret", http.StatusUnauthorized)

		return
	}

	entry, err := h.db.GetChatsByChatUserID(user.ID)
	if err != nil {
		log.Printf("failed to get chat: %v", err)
		util.SendResponse(w, nil, "failed to get chat", http.StatusInternalServerError)

		return
	}

	history := util.ConvertToChatMessage(entry)

	chatHistory := append(history, openai.ChatMessage{
		Role:    openai.ROLE_USER,
		Content: "That is the end of the mock interview, thank you, please provide your feedbacks on my strength and which area to improve, and whether you are confident that I fits the role.",
	})

	answerText, err := util.GenerateText(h.ai, chatHistory)
	if err != nil {
		log.Printf("failed to get chat completion: %v", err)
		util.SendResponse(w, nil, "failed to get chat completion", http.StatusInternalServerError)

		return
	}

	answerAudio, err := util.GenerateSpeech(h.ai, user.Language, answerText)
	if err != nil {
		log.Printf("failed to generate speech: %v", err)
		util.SendResponse(w, nil, "failed to generate speech", http.StatusInternalServerError)

		return
	}

	tx, err := h.db.BeginTx()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		util.SendResponse(w, nil, "failed to create new chat", http.StatusInternalServerError)

		return
	}
	defer tx.Rollback()

	if _, err := h.db.CreateChat(tx, userID, string(openai.ROLE_ASSISTANT), answerText, answerAudio); err != nil {
		log.Printf("failed to create chat: %v", err)
		util.SendResponse(w, nil, "failed to create chat", http.StatusInternalServerError)

		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("failed to commit transaction: %v", err)
		util.SendResponse(w, nil, "failed to create new chat", http.StatusInternalServerError)

		return
	}

	response := model.AnswerChatResponse{
		Language: h.ai.GetCode(user.Language),
		Answer: model.Chat{
			Text:  answerText,
			Audio: answerAudio,
		},
	}

	util.SendResponse(w, response, "success", http.StatusOK)
}
