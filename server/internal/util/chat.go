package util

import (
	"github.com/madeindra/mock-interview/server/internal/data"
	"github.com/madeindra/mock-interview/server/internal/openai"
)

func ConvertToChatMessage(entries []data.Entry) []openai.ChatMessage {
	var messages []openai.ChatMessage
	for _, entry := range entries {
		messages = append(messages, openai.ChatMessage{
			Role:    openai.Role(entry.Role),
			Content: entry.Text,
		})
	}
	return messages
}
