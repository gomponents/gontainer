package template

import (
	"fmt"
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/imports"
	"testing"

	"github.com/gomponents/gontainer-helpers/caller"
	"github.com/stretchr/testify/assert"
)

func Test_createDefaultFunctions(t *testing.T) {
	fncs := createDefaultFunctions(mockAliases{"myAlias"})

	scenarios := []struct {
		fn     string
		input  []interface{}
		output string
		panic  bool
	}{
		{
			fn:     "export",
			input:  []interface{}{5},
			output: "5",
		},
		{
			fn:    "export",
			input: []interface{}{struct{}{}},
			panic: true, // parameter of type `struct {}` is not supported
		},
		{
			fn:     "importAlias",
			input:  []interface{}{"my/package.name"},
			output: "myAlias",
		},
		{
			fn:     "callerAlias",
			output: "myAlias",
		},
		{
			fn:     "containerAlias",
			output: "myAlias",
		},
		{
			fn:     "setterAlias",
			output: "myAlias",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			defer func() {
				if s.panic {
					assert.NotNil(t, recover())
					return
				}

				assert.Nil(t, recover())
			}()

			o := caller.MustCall(fncs[s.fn], s.input...)
			assert.Equal(t, s.output, o[0])
		})
	}
}

func TestSimpleBuilder_Build(t *testing.T) {
	originalHead := templateHead
	originalBody := templateBody

	defer func() {
		templateHead = originalHead
		templateBody = originalBody
	}()

	templateHead = "imports\n(...)\n"
	templateBody = "container(...)"

	o, err := NewSimpleBuilder(mockAliases{alias: "alias"}, mockCollection{}).Build(compiled.DTO{})
	assert.NoError(t, err)
	assert.Equal(t, templateHead+templateBody, o)
}

type mockCollection struct {
	imports []imports.Import
}

func (m mockCollection) GetImports() []imports.Import {
	return m.imports
}

type mockAliases struct {
	alias string
}

func (m mockAliases) GetAlias(string) string {
	return m.alias
}
