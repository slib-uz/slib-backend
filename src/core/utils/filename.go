package utils

import (
	"mime"
	"path/filepath"
	"regexp"
	"strings"
)

// var invalidFilenameChars = regexp.MustCompile(`[^a-zA-Z0-9._-]`)
var invalidFilenameChars = regexp.MustCompile(`[^a-zA-Z0-9._\-\p{Cyrillic}]`)

func SanitizeFilename(filename string) string {
	base := filepath.Base(filename)
	clean := invalidFilenameChars.ReplaceAllString(base, "_")
	return clean
}

func MakeFileName(filename string, suffix string) string {

	name := filepath.Base(filename)

	ext := filepath.Ext(name)
	onlyName := strings.TrimSuffix(name, ext)

	return onlyName + suffix + ext
}

func GetContentType(filePath string) string {
	ext := filepath.Ext(filePath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		return "application/octet-stream"
	}
	return contentType
}
