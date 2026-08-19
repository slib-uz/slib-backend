package utils

import (
	"regexp"
	"strings"
)

func MultiLangValidator(s map[string]string) *string {
	requiredLanguages := []string{"uz", "ru", "en"}

	for _, lang := range requiredLanguages {
		if _, exists := s[lang]; !exists {
			return &lang
		}
	}
	return nil
}

var (
	issnRegex = regexp.MustCompile(`^\d{4}-\d{3}[\dX]$`)
)

func IsValidISSN(issn string) bool {
	if !issnRegex.MatchString(issn) {
		return false
	}

	digits := strings.ReplaceAll(issn, "-", "")

	sum := 0
	for i := range 7 {
		sum += int(digits[i]-'0') * (8 - i)
	}

	remainder := sum % 11
	check := (11 - remainder) % 11

	expected := byte('0' + check)
	if check == 10 {
		expected = 'X'
	}

	return digits[7] == expected
}
