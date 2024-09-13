package config

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

func GetLanguage(code string) string {
	return string(CodeToLanguage[code])
}

func GetCode(lang string) string {
	return LanguageToCode[Language(lang)]
}
