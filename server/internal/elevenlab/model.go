package elevenlab

type TTSRequest struct {
	Text         string       `json:"text"`
	ModelID      string       `json:"model_id"`
	VoiceSetting VoiceSetting `json:"voice_settings"`
}

type VoiceSetting struct {
	Stability       float32 `json:"stability"`
	SimilarityBoost float32 `json:"similarity_boost"`
}
