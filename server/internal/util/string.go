package util

import (
	"fmt"
	"regexp"
	"strings"
)

func SanitizeString(text string) string {
	reStrong := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	text = reStrong.ReplaceAllString(text, "$1")

	reItalic := regexp.MustCompile(`\*([^*]+)\*`)
	text = reItalic.ReplaceAllString(text, "$1")

	reLink := regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
	text = reLink.ReplaceAllString(text, "$1")

	reBullet := regexp.MustCompile(`\n- `)
	text = reBullet.ReplaceAllString(text, ", ")

	text = strings.Replace(text, "\n", " ", -1)

	return text
}

func SanitizeSSML(input string) (string, error) {
	if !strings.Contains(input, "<speak") || !strings.Contains(input, "</speak>") {
		return "", fmt.Errorf("input string does not contain both <speak and </speak> tags")
	}

	startIndex := strings.Index(input, "<speak")
	if startIndex == -1 {
		return "", fmt.Errorf("could not find opening <speak tag")
	}

	endIndex := strings.LastIndex(input, "</speak>")
	if endIndex == -1 {
		return "", fmt.Errorf("could not find closing </speak> tag")
	}

	return input[startIndex : endIndex+8], nil
}

// TODO - strip all SSML tag and compare whether both are identical
func ValidateIdentical(original, ssml string) error {
	return nil
}
