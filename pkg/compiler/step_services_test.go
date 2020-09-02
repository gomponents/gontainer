package compiler

import "testing"

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
