package tmpl

import (
	"embed"
	"fmt"
	"text/template"
)

var (
	//go:embed *.tmpl
	FS embed.FS
)

const (
	TemplateEmail        = "email"
	TemplateEmailSubject = "subject"
	TemplateEmailBody    = "body"
)

// Load parses templates in the current directory.
func Load() (*template.Template, error) {
	tmpl, err := template.ParseFS(FS, "*")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	return tmpl, nil
}
