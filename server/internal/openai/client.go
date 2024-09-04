package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type Client interface {
	IsKeyValid() (bool, error)
	Status() (Status, error)
	Chat([]ChatMessage) (ChatResponse, error)
	TextToSpeech(string) (io.ReadCloser, error)
	Transcribe(io.ReadCloser, string, string) (TranscriptResponse, error)

	GetDefaultTranscriptLanguage() string
	GetLanguage(code string) string
	GetCode(lang string) string
	IsSpeechAvailable(string) bool
}

type OpenAI struct {
	APIKey             string
	BaseURL            string
	ChatModel          string
	TranscriptModel    string
	TranscriptLanguage string
	TTSModel           string
	TTSVoice           string
}

const (
	baseURL            = "https://api.openai.com/v1"
	statusURL          = "https://status.openai.com/api/v2"
	chatModel          = "gpt-4o"
	transcriptModel    = "whisper-1"
	transcriptLanguage = "en"
	ttsModel           = "tts-1"
	ttsVoice           = "nova"
)

var supportedTranscriptLanguages = map[Language]struct{}{
	LANGUAGE_ENGLISH: {},
}

func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{
		APIKey:             apiKey,
		BaseURL:            baseURL,
		ChatModel:          chatModel,
		TranscriptModel:    transcriptModel,
		TTSModel:           ttsModel,
		TTSVoice:           ttsVoice,
		TranscriptLanguage: transcriptLanguage,
	}
}

func (c *OpenAI) IsKeyValid() (bool, error) {
	url, err := url.JoinPath(c.BaseURL, "/models")
	if err != nil {
		return false, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}

func (c *OpenAI) Status() (Status, error) {
	url, err := url.JoinPath(statusURL, "/components.json")
	if err != nil {
		return STATUS_UNKNOWN, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return STATUS_UNKNOWN, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return STATUS_UNKNOWN, err
	}

	if resp.StatusCode != http.StatusOK {
		return STATUS_UNKNOWN, nil
	}

	var statusResp ComponentStatusResponse
	err = unmarshalJSONResponse(resp, &statusResp)
	if err != nil {
		return STATUS_UNKNOWN, err
	}

	for _, component := range statusResp.Components {
		if component.Name == "API" {
			switch component.Status {
			case "operational":
				return STATUS_OPERATIONAL, nil
			case "degraded_performance":
				return STATUS_DEGRADED_PERFORMANCE, nil
			case "partial_outage":
				return STATUS_PARTIAL_OUTAGE, nil
			case "major_outage":
				return STATUS_MAJOR_OUTAGE, nil
			}
		}
	}

	return STATUS_UNKNOWN, nil
}

func (c *OpenAI) Chat(messages []ChatMessage) (ChatResponse, error) {
	url, err := url.JoinPath(c.BaseURL, "/chat/completions")
	if err != nil {
		return ChatResponse{}, err
	}

	chatReq := ChatRequest{
		Model:    c.ChatModel,
		Messages: messages,
	}

	body, err := json.Marshal(chatReq)
	if err != nil {
		return ChatResponse{}, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return ChatResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ChatResponse{}, err
	}

	var chatResp ChatResponse
	err = unmarshalJSONResponse(resp, &chatResp)
	if err != nil {
		return ChatResponse{}, err
	}

	return chatResp, nil
}

func (c *OpenAI) TextToSpeech(input string) (io.ReadCloser, error) {
	url, err := url.JoinPath(c.BaseURL, "/audio/speech")
	if err != nil {
		return nil, err
	}

	ttsReq := TTSRequest{
		Model: c.TTSModel,
		Voice: c.TTSVoice,
		Input: input,
	}

	body, err := json.Marshal(ttsReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, err := getResponseBody(resp)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *OpenAI) Transcribe(file io.ReadCloser, filename, language string) (TranscriptResponse, error) {
	if file == nil {
		return TranscriptResponse{}, fmt.Errorf("audio is nil")
	}
	defer file.Close()

	url, err := url.JoinPath(c.BaseURL, "/audio/transcriptions")
	if err != nil {
		return TranscriptResponse{}, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return TranscriptResponse{}, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return TranscriptResponse{}, err
	}

	err = writer.WriteField("model", c.TranscriptModel)
	if err != nil {
		return TranscriptResponse{}, err
	}

	transcriptLanguage := c.TranscriptLanguage
	if language != "" {
		transcriptLanguage = language
	}

	err = writer.WriteField("language", transcriptLanguage)
	if err != nil {
		return TranscriptResponse{}, err
	}

	err = writer.Close()
	if err != nil {
		return TranscriptResponse{}, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, body)
	if err != nil {
		return TranscriptResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return TranscriptResponse{}, err
	}

	if resp == nil || resp.Body == nil {
		return TranscriptResponse{}, fmt.Errorf("response is nil")
	}

	var transcriptResp TranscriptResponse
	err = unmarshalJSONResponse(resp, &transcriptResp)
	if err != nil {
		return TranscriptResponse{}, err
	}

	return transcriptResp, nil
}

func (c *OpenAI) GetDefaultTranscriptLanguage() string {
	return string(Language(c.TranscriptLanguage))
}

func (c *OpenAI) GetLanguage(code string) string {
	return string(CodeToLanguage[code])
}

func (c *OpenAI) GetCode(lang string) string {
	return LanguageToCode[Language(lang)]
}

func (c *OpenAI) IsSpeechAvailable(lang string) bool {
	_, ok := supportedTranscriptLanguages[Language(lang)]
	return ok
}

func getResponseBody(resp *http.Response) (io.ReadCloser, error) {
	if resp == nil || resp.Body == nil {
		return nil, fmt.Errorf("response is nil")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func unmarshalJSONResponse(resp *http.Response, v interface{}) error {
	respBody, err := getResponseBody(resp)
	if err != nil {
		return err
	}
	if respBody == nil {
		return fmt.Errorf("response body is nil")
	}
	defer respBody.Close()

	respByte, err := io.ReadAll(respBody)
	if err != nil {
		return err
	}

	return json.Unmarshal(respByte, v)
}
