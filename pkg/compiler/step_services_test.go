package compiler

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/stretchr/testify/assert"
)

func TestStepServices_handleServiceType(t *testing.T) {
	scenarios := []inputOutputScenario{
		{
			input:  "MyStruct",
			output: "MyStruct",
		},
		{
			input:  "my/import.MyStruct",
			output: "alias.MyStruct",
		},
		{
			input:  `"my/import".MyStruct`,
			output: "alias.MyStruct",
		},
		{
			input:  "",
			output: "",
		},
	}

	doTestInputOutput(
		t,
		StepServices{
			aliases: mockImports{alias: "alias"},
		}.handleServiceType,
		scenarios...,
	)
}

func TestStepServices_handleServiceValue(t *testing.T) {
	scenarios := []inputOutputScenario{
		{
			input:  "MyValue",
			output: "MyValue",
		},
		{
			input:  `".".MyValue.Foo`,
			output: "MyValue.Foo",
		},
		{
			input:  "&MyValue",
			output: "&MyValue",
		},
		{
			input:  `&".".MyValue`,
			output: "&MyValue",
		},
		{
			input:  "my/import.MyValue",
			output: "alias.MyValue",
		},
		{
			input:  "&my/import.MyValue",
			output: "&alias.MyValue",
		},
		{
			input:  `"my/import".MyValue`,
			output: "alias.MyValue",
		},
		{
			input:  `&"my/import".MyValue`,
			output: "&alias.MyValue",
		},
		{
			input:  `my/import.MyVar.SomeField`, // compiler doesn't know whether `my/import` or `my/import.MyVar` is the import path
			output: "alias.SomeField",
		},
		{
			input:  `"my/import".MyVar.SomeField`, // surrounding import by `"` makes it explicit
			output: "alias.MyVar.SomeField",
		},
		{
			input:  `"github.com/gontainer/gontainer/pkg/dto/input".GlobalDTO.Meta.Functions`,
			output: `alias.GlobalDTO.Meta.Functions`,
		},
		{
			input:  `&"my/import".MyStruct{}`,
			output: "&alias.MyStruct{}",
		},
		{
			input:  `&my/import.MyStruct{}`,
			output: "&alias.MyStruct{}",
		},
		{
			input:  `&MyStruct{}`,
			output: "&MyStruct{}",
		},
		{
			input:  `MyStruct{}`,
			output: "MyStruct{}",
		},
		{
			input:  "",
			output: "",
		},
	}

	doTestInputOutput(
		t,
		StepServices{
			aliases: mockImports{alias: "alias"},
		}.handleServiceValue,
		scenarios...,
	)
}

func TestStepServices_handleServiceConstructor(t *testing.T) {
	scenarios := []inputOutputScenario{
		{
			input:  "my/import.NewFoo",
			output: "alias.NewFoo",
		},
		{
			input:  `"my/import".NewBar`,
			output: "alias.NewBar",
		},
		{
			input:  "NewFoo",
			output: "NewFoo",
		},
		{
			input:  "",
			output: "",
		},
	}

	doTestInputOutput(
		t,
		StepServices{
			aliases: mockImports{alias: "alias"},
		}.handleServiceConstructor,
		scenarios...,
	)
}

func TestStepServices_Do(t *testing.T) {
	t.Run("Sort services", func(t *testing.T) {
		scenarios := []struct {
			input  map[string]input.Service
			output []compiled.Service
		}{
			{
				input: map[string]input.Service{
					"logger":     {Todo: true},
					"db":         {Todo: true},
					"httpClient": {Todo: true},
				},
				output: []compiled.Service{
					{Name: "db", Todo: true},
					{Name: "httpClient", Todo: true},
					{Name: "logger", Todo: true},
				},
			},
			{
				input:  nil,
				output: nil,
			},
		}

		for i, s := range scenarios {
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				i := input.DTO{Services: s.input}
				o := compiled.DTO{}
				err := StepServices{}.Do(&i, &o)
				assert.NoError(t, err)
				assert.Equal(t, s.output, o.Services)
			})
		}
	})
}

func TestStepServices_handleServiceTags(t *testing.T) {
	scenarios := []inputOutputScenario{
		{
			input: []input.Tag{
				{Name: "storage", Priority: 200},
				{Name: "db", Priority: 100},
			},
			output: []compiled.Tag{
				{Name: "storage", Priority: 200},
				{Name: "db", Priority: 100},
			},
		},
		{
			input:  ([]input.Tag)(nil),
			output: ([]compiled.Tag)(nil),
		},
	}

	doTestInputOutput(
		t,
		StepServices{}.handleServiceTags,
		scenarios...,
	)
}

type mockImports struct {
	alias string
	error error
}

func (m mockImports) GetAlias(string) string {
	return m.alias
}

func (m mockImports) RegisterPrefix(shortcut string, path string) error {
	return m.error
}
