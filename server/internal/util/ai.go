package util

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/madeindra/mock-interview/server/internal/openai"
)

func GetChatAssets(ai openai.Client, role string, skills []string, language string) (string, string, error) {
	if ai == nil {
		return "", "", fmt.Errorf("unsupported client")
	}

	systempPrompt, err := openai.GetSystemPrompt(role, skills, language)
	if err != nil {
		return "", "", err
	}

	initialChat, err := openai.GetInitialChat(role, language)
	if err != nil {
		return "", "", err
	}

	return systempPrompt, initialChat, nil
}

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

	if !ai.IsSpeechAvailable(language) {
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

func GenerateSSML(ai openai.Client, text string) (string, error) {
	if ai == nil {
		return "", fmt.Errorf("unsupported client")
	}

	generated, err := ai.SSML(text)
	if err != nil {
		return "", err
	}

	if generated.SSML == "" {
		return "", fmt.Errorf("empty ssml response")
	}

	sanitized, err := SanitizeSSML(generated.SSML)
	if err != nil {
		return "", nil // quietly ignore improper formatted response
	}

	return sanitized, nil
}
