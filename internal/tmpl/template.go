package tmpl

import (
	"embed"
	"fmt"
	"html/template"
	"io"
)

type Template struct {
	templates *template.Template
}

func NewTemplateRenderer(fs embed.FS, patterns ...string) (Template, error) {
	funcMap := template.FuncMap{}

	t, err := template.New("").Funcs(funcMap).ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("unable to parse templates: %w", err)
	}

	return Template{templates: t}, nil
}

func (t Template) Render(w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
