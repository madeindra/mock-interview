package elevenlab

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client interface {
	TextToSpeech(string) (io.ReadCloser, error)
}

type ElevenLab struct {
	APIKey   string
	BaseURL  string
	TTSModel string
	TTSVoice string
}

const (
	baseURL  = "https://api.elevenlabs.io/v1"
	ttsModel = "eleven_multilingual_v2"
	ttsVoice = "cgSgspJ2msm6clMCkdW9"
)

var defaultVoiceSetting = VoiceSetting{
	Stability:       0.5,
	SimilarityBoost: 0.75,
}

func NewElevenLab(apiKey string) *ElevenLab {
	return &ElevenLab{
		APIKey:   apiKey,
		BaseURL:  baseURL,
		TTSModel: ttsModel,
		TTSVoice: ttsVoice,
	}
}

func (c *ElevenLab) TextToSpeech(input string) (io.ReadCloser, error) {
	url, err := url.JoinPath(c.BaseURL, "text-to-speech", c.TTSVoice)
	if err != nil {
		return nil, err
	}

	ttsReq := TTSRequest{
		Text:         input,
		ModelID:      c.TTSModel,
		VoiceSetting: defaultVoiceSetting,
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

func getResponseBody(resp *http.Response) (io.ReadCloser, error) {
	if resp == nil || resp.Body == nil {
		return nil, fmt.Errorf("response is nil")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
