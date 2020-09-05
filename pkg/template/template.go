package template

import (
	"bytes"
	"text/template"

	"github.com/gomponents/gontainer-helpers/exporters"
	"github.com/gomponents/gontainer/pkg/consts"
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/imports"
)

//go:generate go run ../../embed-file/main.go head.tmpl template templateHead tmpl_head.go
//go:generate go run ../../embed-file/main.go body.tmpl template templateBody tmpl_body.go

type Builder struct {
	aliases    imports.Aliases
	collection imports.Collection
}

func NewBuilder(aliases imports.Aliases, collection imports.Collection) *Builder {
	return &Builder{aliases: aliases, collection: collection}
}

func (b Builder) Build(i compiled.DTO) (string, error) {
	data := map[string]interface{}{
		"ImportCollection": b.collection,
		"Input":            i,
	}

	funcs := createDefaultFunctions(b.aliases)

	var (
		body, head string
		err        error
	)

	tplBody := tpl{
		name:  "body",
		body:  templateBody,
		vars:  data,
		funcs: funcs,
	}

	tplHead := tpl{
		name:  "head",
		body:  templateHead,
		vars:  data,
		funcs: funcs,
	}

	if body, err = tplBody.exec(); err != nil {
		return "", err
	}

	if head, err = tplHead.exec(); err != nil {
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

type tpl struct {
	name  string
	body  string
	vars  map[string]interface{}
	funcs template.FuncMap
}

func (t tpl) exec() (string, error) {
	tpl, newErr := template.New("gontainer_" + t.name).Funcs(t.funcs).Parse(t.body)
	if newErr != nil {
		return "", newErr
	}
	var b bytes.Buffer
	tplErr := tpl.Execute(&b, t.vars)
	return b.String(), tplErr
}
