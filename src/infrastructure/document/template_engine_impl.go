package document

import (
	"bytes"
	"html/template"

	"slib.uz/src/core/domain/ports/document"
)

type TemplateEngineImpl struct {
}

// @inject
func NewTemplateEngine() document.TemplateEngine {
	return &TemplateEngineImpl{}
}

func (this *TemplateEngineImpl) Render(templatePath string, data any) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil

}
