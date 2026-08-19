package entity

import (
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"
)

type ArticleDownloadEntity struct {
	PresignedURL string `json:"presigned_url"`
	BinaryFile   string `json:"binary_file"`
}

func NewArticleDownloadEntity(presignedURL string, binaryFile []byte) *ArticleDownloadEntity {
	// Convert binary to base64
	base64Data := base64.StdEncoding.EncodeToString(binaryFile)

	// Detect MIME type from URL (simple approach)
	mimeType := detectMimeType(presignedURL)

	// Format: data:<mimeType>;base64,<base64Data>
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data)

	return &ArticleDownloadEntity{
		PresignedURL: presignedURL,
		BinaryFile:   dataURL,
	}
}

// detectMimeType detects MIME type from file extension
func detectMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	mimeTypes := map[string]string{
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".txt":  "text/plain",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
	}

	if mimeType, ok := mimeTypes[ext]; ok {
		return mimeType
	}

	// Default to application/octet-stream for unknown types
	return "application/pdf"
}
