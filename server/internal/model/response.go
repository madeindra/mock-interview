package model

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type StartChatResponse struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`

	Chat
}

type AnswerChatResponse struct {
	Prompt Chat `json:"prompt,omitempty"`
	Answer Chat `json:"answer,omitempty"`
}

type StatusResponse struct {
	Server    bool    `json:"backend"`
	API       *bool   `json:"api"`
	ApiStatus *string `json:"apiStatus"`
	Key       bool    `json:"key"`
}
