package model

type StartChatRequest struct {
	Role   string   `json:"role"`
	Skills []string `json:"skills"`
}
