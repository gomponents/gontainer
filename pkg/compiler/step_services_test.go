package compiler

import (
	"testing"
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
			input:  "my/import.MyValue",
			output: "alias.MyValue",
		},
		{
			input:  `"my/import".MyValue`,
			output: "alias.MyValue",
		},
		{
			input:  `"my/import".MyStruct{}.MyMethod`,
			output: "alias.MyStruct{}.MyMethod",
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
	}

	doTestInputOutput(
		t,
		StepServices{
			aliases: mockImports{alias: "alias"},
		}.handleServiceConstructor,
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
