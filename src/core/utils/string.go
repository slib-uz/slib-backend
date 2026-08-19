package utils

import "unicode/utf8"

func Substring(str string, from, to int) string {
	runes := []rune(str)
	length := len(runes)

	if from < 0 {
		from = 0
	}
	if to > length {
		to = length
	}
	if from > to {
		return ""
	}
	return string(runes[from:to])

}

func TruncateUTF8(str string, maxLen int) string {
	if utf8.RuneCountInString(str) <= maxLen {
		return str
	}
	runes := []rune(str)
	return string(runes[:maxLen])
}
