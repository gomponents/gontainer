package template

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/gomponents/gontainer-helpers/exporters"
	"github.com/gomponents/gontainer/pkg/consts"
	"github.com/gomponents/gontainer/pkg/dto/compiled"
)

//go:generate go run ../../templater/main.go head.tmpl template TemplateHead tmpl_head.go
//go:generate go run ../../templater/main.go body.tmpl template TemplateBody tmpl_body.go

type Imports interface {
	GetAlias(string) string
}

type SimpleBuilder struct {
	imports Imports
}

func NewSimpleBuilder(imports Imports) *SimpleBuilder {
	return &SimpleBuilder{imports: imports}
}

func (s SimpleBuilder) Build(i compiled.DTO) (string, error) {
	data := map[string]interface{}{
		"Imports": s.imports,
		"Input":   i,
	}

	fncs := template.FuncMap{
		"export": func(input interface{}) string {
			r, err := exporters.NewDefaultExporter().Export(input)
			if err != nil {
				panic(err)
			}
			return r
		},
		"importAlias": func(i string) string {
			return s.imports.GetAlias(i)
		},
		"replace": func(input, from, to string) string {
			return strings.Replace(input, from, to, -1)
		},
		"callerAlias": func() string {
			return s.imports.GetAlias(consts.GontainerHelperPath + "/caller")
		},
		"containerAlias": func() string {
			return s.imports.GetAlias(consts.GontainerHelperPath + "/container")
		},
		"setterAlias": func() string {
			return s.imports.GetAlias(consts.GontainerHelperPath + "/setter")
		},
	}

	exec := func(name string, tplBody string) (string, error) {
		tpl, newErr := template.New("gontainer_" + name).Funcs(fncs).Parse(tplBody)
		if newErr != nil {
			return "", newErr
		}
		var b bytes.Buffer
		tplErr := tpl.Execute(&b, data)
		return b.String(), tplErr
	}

	var (
		body, head string
		err        error
	)

	if body, err = exec("body", TemplateBody); err != nil {
		return "", err
	}

	if head, err = exec("head", TemplateHead); err != nil {
		return "", err
	}

	return head + body, nil
}
