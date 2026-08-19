package pdf_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"slib.uz/src/infrastructure/pdf"
)

func TestExtractReturnsErrorForNonPdf(t *testing.T) {
	_, err := pdf.NewTextExtractor().Extract([]byte("not a pdf"))
	if err == nil {
		t.Fatal("kutilgan xato")
	}
}

func TestExtractReadsSimpleTextLayer(t *testing.T) {
	text, err := pdf.NewTextExtractor().Extract(helloPDF())
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if !strings.Contains(text, "Hello") {
		t.Fatalf("text=%q", text)
	}
}

func helloPDF() []byte {
	stream := "BT /F1 12 Tf 10 100 Td (Hello) Tj ET"
	objects := []string{
		"1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n",
		"2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n",
		"3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 200 200] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>\nendobj\n",
		fmt.Sprintf("4 0 obj\n<< /Length %d >>\nstream\n%s\nendstream\nendobj\n", len(stream), stream),
		"5 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n",
	}

	var body bytes.Buffer
	body.WriteString("%PDF-1.4\n")
	offsets := make([]int, len(objects)+1)
	for i, obj := range objects {
		offsets[i+1] = body.Len()
		body.WriteString(obj)
	}
	startxref := body.Len()
	fmt.Fprintf(&body, "xref\n0 %d\n", len(objects)+1)
	body.WriteString("0000000000 65535 f \n")
	for i := 1; i <= len(objects); i++ {
		fmt.Fprintf(&body, "%010d 00000 n \n", offsets[i])
	}
	fmt.Fprintf(&body, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objects)+1, startxref)
	return body.Bytes()
}
