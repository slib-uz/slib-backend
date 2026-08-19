package utils

import (
	"errors"
	"strings"
	"unicode/utf8"
)

const (
	maxTagNameLength      = 32
	maxLanguageCodeLength = 5
)

var (
	ErrTagNameTooLong      = errors.New("tag name must be at most 32 characters")
	ErrLanguageCodeTooLong = errors.New("language code must be at most 5 characters")
)

type TagNameLang struct {
	Lang string
	Name string
}

func NormalizeTagNamesByLang(in map[string][]string) ([]TagNameLang, error) {
	out := make([]TagNameLang, 0)
	seen := make(map[string]bool)
	for lang, names := range in {
		lang = strings.ToLower(strings.TrimSpace(lang))
		if lang == "" {
			continue
		}
		if utf8.RuneCountInString(lang) > maxLanguageCodeLength {
			return nil, ErrLanguageCodeTooLong
		}
		for _, name := range names {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			if utf8.RuneCountInString(name) > maxTagNameLength {
				return nil, ErrTagNameTooLong
			}
			key := lang + "\x00" + strings.ToLower(name)
			if seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, TagNameLang{Lang: lang, Name: name})
		}
	}
	return out, nil
}
