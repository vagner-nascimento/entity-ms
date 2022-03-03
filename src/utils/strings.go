package utils

import (
	"unicode"
	"unicode/utf8"
)

func LowerFirst(text string) string {
	if len(text) == 0 {
		return ""
	}

	first, index := utf8.DecodeRuneInString(text)

	return string(unicode.ToLower(first)) + text[index:]
}
