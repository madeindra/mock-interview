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
	//go:embed templates/chat.txt
	initalChat string

	//go:embed templates/system.txt
	systemPrompt string
)

func GetSystemPrompt(roleName string, skills []string) (string, error) {
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

func GetInitialChat(roleName string) (string, error) {
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
