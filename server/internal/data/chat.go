package data

import "github.com/madeindra/mock-interview/server/internal/openai"

type ChatEntry struct {
	ID      string `bson:"_id"`
	Secret  string
	History []openai.ChatMessage
}
