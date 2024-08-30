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
	Chat([]ChatMessage) (ChatResponse, error)
	TextToSpeech(string) (io.ReadCloser, error)
	Transcribe(io.ReadCloser, string) (TranscriptResponse, error)
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
	chatModel          = "gpt-4o"
	transcriptModel    = "whisper-1"
	transcriptLanguage = "en"
	ttsModel           = "tts-1"
	ttsVoice           = "nova"
)

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

func (c *OpenAI) Transcribe(file io.ReadCloser, filename string) (TranscriptResponse, error) {
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

	err = writer.WriteField("language", c.TranscriptLanguage)
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
