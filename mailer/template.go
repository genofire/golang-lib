package mailer

import (
	"bytes"
	hTemplate "html/template"
	"text/template"
)

// Renderer for easy template of TXT or HTML
type Renderer struct {
	tmpl  *template.Template
	hTmpl *hTemplate.Template
}

// TemplateTXT - create template render
func TemplateTXT(temp string) *Renderer {
	tmpl, err := template.New("").Parse(temp)
	if err != nil {
		return nil
	}
	return &Renderer{
		tmpl: tmpl,
	}
}

// TemplateHTML - create template render for html
func TemplateHTML(temp string) *Renderer {
	tmpl, err := hTemplate.New("").Parse(temp)
	if err != nil {
		return nil
	}
	return &Renderer{
		hTmpl: tmpl,
	}
}

// Render template
func (r *Renderer) Render(data interface{}) string {
	var buf bytes.Buffer
	if r.hTmpl != nil {
		if err := r.hTmpl.Execute(&buf, data); err != nil {
			return ""
		}
	} else {
		if err := r.tmpl.Execute(&buf, data); err != nil {
			return ""
		}
	}
	return string(buf.Bytes())
}
