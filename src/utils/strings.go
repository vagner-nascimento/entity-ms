package utils

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// return the sent text with the first letter in lower case
func LowerFirst(text string) (res string) {
	if len(text) > 0 {
		first, index := utf8.DecodeRuneInString(text)
		res = string(unicode.ToLower(first)) + text[index:]
	}

	return
}

// Returns the first string found between brackets symbol "[]" OR empty string if none was found
func StringBetweenBrackets(s string) (res string) {
	match := strings.Split(s, "[")
	if len(match) > 1 {
		res = strings.Split(match[1], "]")[0]
	}

	return
}
