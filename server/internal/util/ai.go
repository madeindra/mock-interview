package util

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/madeindra/mock-interview/server/internal/openai"
)

func TranscribeSpeech(ai openai.Client, file io.ReadCloser, filename, language string) (string, error) {
	if ai == nil {
		return "", fmt.Errorf("unsupported client")
	}

	transcript, err := ai.Transcribe(file, filename, language)
	if err != nil {
		return "", err
	}

	if transcript.Text == "" {
		return "", fmt.Errorf("empty transcript")
	}

	return transcript.Text, nil
}

func GenerateText(ai openai.Client, entries []openai.ChatMessage) (string, error) {
	if ai == nil {
		return "", fmt.Errorf("unsupported client")
	}

	chatCompletion, err := ai.Chat(entries)
	if err != nil {
		return "", err
	}

	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("empty chat response")
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func GenerateSpeech(ai openai.Client, language, text string) (string, error) {
	if ai == nil {
		return "", fmt.Errorf("unsupported client")
	}

	if ai.IsSpeechAvailable(language) {
		return "", nil // quietly ignore unsupported languages
	}

	speechInput := SanitizeString(text)

	speech, err := ai.TextToSpeech(speechInput)
	if err != nil {
		return "", err
	}

	speechByte, err := io.ReadAll(speech)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(speechByte), nil
}
