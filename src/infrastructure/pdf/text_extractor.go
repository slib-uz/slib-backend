package pdf

import (
	"bytes"
	"fmt"
	"strings"

	ledongthucpdf "github.com/ledongthuc/pdf"
	domainpdf "slib.uz/src/core/domain/ports/pdf"
)

type TextExtractorImpl struct{}

// @inject
func NewTextExtractor() domainpdf.TextExtractor {
	return &TextExtractorImpl{}
}

func (this *TextExtractorImpl) Extract(data []byte) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("empty pdf")
	}

	reader, err := ledongthucpdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}

	var b strings.Builder
	numPage := reader.NumPage()
	for i := 1; i <= numPage; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(text) == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(text)
	}

	return b.String(), nil
}
