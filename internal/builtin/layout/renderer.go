package layout

import (
	"bytes"
	"io"
	"text/template"

	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/pkg/errors"
)

type (
	// RendererInterface define common interface for rendering templates.
	RendererInterface interface {
		Render(source string, values generator.RunValues) ([]byte, error)
	}

	// TemplateInterface define common interface for parsing and executing text/html templates.
	TemplateInterface interface {
		Parse(string) (*template.Template, error)
		Execute(io.Writer, interface{}) error
	}

	// Renderer is built-in RendererInterface implementation.
	Renderer struct {
		box      BoxInterface
		template TemplateInterface
	}
)

var (
	_ RendererInterface = &Renderer{}
)

// NewTemplate return built-in go *template.Template.
func NewTemplate(name string, funcMap template.FuncMap) *template.Template {
	t := template.New(name)
	if funcMap != nil {
		t.Funcs(funcMap)
	}
	return t
}

// NewRenderer return implementation of RendererInterface.
func NewRenderer(box BoxInterface, ti TemplateInterface) RendererInterface {
	return &Renderer{
		box:      box,
		template: ti,
	}
}

/*
Render implements RendererInterface
Steps for render:
- try to find template in embed box or resolve in in os filesystem
- open file
- get template file content
- parse template file content
- execute template
*/
func (r *Renderer) Render(source string, values generator.RunValues) ([]byte, error) {
	content, err := r.box.ReadFile(source)
	if err != nil {
		return nil, errors.Wrapf(err, `failed to get source file "%s" content`, source)
	}
	t, err := r.template.Parse(string(content))
	if err != nil {
		return nil, errors.Wrapf(err, `failed to parse source file "%s"`, source)
	}

	b := new(bytes.Buffer)
	if err = t.Execute(b, values); err != nil {
		return nil, errors.Wrapf(err, `failed to render source file "%s"`, source)
	}
	return b.Bytes(), nil
}
