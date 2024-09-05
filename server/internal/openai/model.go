package openai

type ChatRequest struct {
	Messages       []ChatMessage   `json:"messages"`
	Model          string          `json:"model"`
	ResponseFormat *map[string]any `json:"responseFormat"`
}

type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

type SSMLResponse struct {
	SSML string `json:"ssml"`
}

type ChatMessage struct {
	Content string `json:"content"`
	Role    Role   `json:"role"`
}

type Role string

const (
	ROLE_SYSTEM    Role = "system"
	ROLE_ASSISTANT Role = "assistant"
	ROLE_USER      Role = "user"
)

type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type TTSRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
	Voice string `json:"voice"`
}

type TranscriptResponse struct {
	Text string `json:"text"`
}

type Status string

const (
	STATUS_OPERATIONAL          Status = "operational"
	STATUS_DEGRADED_PERFORMANCE Status = "degraded_performance"
	STATUS_PARTIAL_OUTAGE       Status = "partial_outage"
	STATUS_MAJOR_OUTAGE         Status = "major_outage"
	STATUS_UNKNOWN              Status = "unknown"
)

type ComponentStatusResponse struct {
	Components []Component `json:"components"`
}

type Component struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
}

type Language string

const (
	LANGUAGE_ENGLISH    Language = "en"
	LANGUAGE_INDONESIAN Language = "id"

	CODE_ENGLISH    = "en-US"
	CODE_INDONESIAN = "id-ID"
)

var CodeToLanguage = map[string]Language{
	CODE_ENGLISH:    LANGUAGE_ENGLISH,
	CODE_INDONESIAN: LANGUAGE_INDONESIAN,
}

var LanguageToCode = map[Language]string{
	LANGUAGE_ENGLISH:    CODE_ENGLISH,
	LANGUAGE_INDONESIAN: CODE_INDONESIAN,
}
