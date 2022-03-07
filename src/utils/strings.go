package utils

import (
	"unicode"
	"unicode/utf8"
)

func LowerFirst(text string) (res string) {
	if len(text) > 0 {
		first, index := utf8.DecodeRuneInString(text)
		res = string(unicode.ToLower(first)) + text[index:]
	}

	return
}
