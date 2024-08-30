package openai

type ChatRequest struct {
	Messages []ChatMessage `json:"messages"`
	Model    string        `json:"model"`
}

type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

type ChatMessage struct {
	Content string `json:"content"`
	Role    Role   `json:"role"`
}

type Role string

const (
	ROLE_SYSTEM    Role = "system"
	ROLE_ASSISTANT Role = "assistant"
	ROLE_USER      Role = "user"
)

type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type TTSRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
	Voice string `json:"voice"`
}

type TranscriptResponse struct {
	Text string `json:"text"`
}
