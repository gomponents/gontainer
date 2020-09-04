package template

import (
	"bytes"
	"text/template"

	"github.com/gomponents/gontainer-helpers/exporters"
	"github.com/gomponents/gontainer/pkg/consts"
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/imports"
)

//go:generate go run ../../templater/main.go head.tmpl template TemplateHead tmpl_head.go
//go:generate go run ../../templater/main.go body.tmpl template TemplateBody tmpl_body.go

type SimpleBuilder struct {
	aliases    imports.Aliases
	collection imports.Collection
}

func NewSimpleBuilder(aliases imports.Aliases, collection imports.Collection) *SimpleBuilder {
	return &SimpleBuilder{aliases: aliases, collection: collection}
}

func (s SimpleBuilder) Build(i compiled.DTO) (string, error) {
	data := map[string]interface{}{
		"ImportCollection": s.collection,
		"Input":            i,
	}

	fncs := createDefaultFunctions(s.aliases)

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

func createDefaultFunctions(a imports.Aliases) template.FuncMap {
	return template.FuncMap{
		"export": func(input interface{}) string {
			r, err := exporters.NewDefaultExporter().Export(input)
			if err != nil {
				panic(err)
			}
			return r
		},
		"importAlias": func(i string) string {
			return a.GetAlias(i)
		},
		"callerAlias": func() string {
			return a.GetAlias(consts.GontainerHelperPath + "/caller")
		},
		"containerAlias": func() string {
			return a.GetAlias(consts.GontainerHelperPath + "/container")
		},
		"setterAlias": func() string {
			return a.GetAlias(consts.GontainerHelperPath + "/setter")
		},
	}
}
