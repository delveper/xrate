package tmpl

import (
	"embed"
	"fmt"
	"text/template"
)

var (
	//go:embed *.tmpl
	fS embed.FS
)

const (
	TemplateEmail        = "email"
	TemplateEmailSubject = "subject"
	TemplateEmailBody    = "body"
)

// Load parses templates in the current directory.
func Load() (*template.Template, error) {
	// TODO: improve flexibility to choose templates.
	tmpl, err := template.ParseFS(fS, "*")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	return tmpl, nil
}
