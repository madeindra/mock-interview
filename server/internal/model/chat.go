package model

type Chat struct {
	Audio string `json:"audio,omitempty"`
	SSML  string `json:"ssml,omitempty"`
	Text  string `json:"text,omitempty"`
}
