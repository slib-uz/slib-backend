package document

type TemplateEngine interface {
	Render(templatePath string, data any) (string, error)
}
