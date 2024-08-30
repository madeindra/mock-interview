package util

import (
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
