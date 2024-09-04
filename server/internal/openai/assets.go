package openai

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

type ChatAsset struct {
	SystemPrompt string
	ChatText     string
	ChatAudio    string
}

var (
	//go:embed templates/chat.en.txt
	initalChatEN string

	//go:embed templates/chat.id.txt
	initalChatID string

	//go:embed templates/system.en.txt
	systemPromptEN string

	//go:embed templates/system.id.txt
	systemPromptID string
)

func GetSystemPrompt(roleName string, skills []string, language string) (string, error) {
	systemPrompt := systemPromptEN
	if language == "id" {
		systemPrompt = systemPromptID
	}

	t, err := template.New("prompt").Parse(systemPrompt)
	if err != nil {
		return "", err
	}

	data := struct {
		Role   string
		Skills string
	}{
		Role:   roleName,
		Skills: strings.Join(skills, ";"),
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetInitialChat(roleName string, language string) (string, error) {
	initalChat := initalChatEN
	if language == "id" {
		initalChat = initalChatID
	}

	t, err := template.New("chat").Parse(initalChat)
	if err != nil {
		return "", err
	}

	data := struct {
		Role string
	}{
		Role: roleName,
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
