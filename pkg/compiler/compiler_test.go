package compiler

import (
	"fmt"
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompiler_handleService(t *testing.T) {
	scenarios := []struct {
		argResolver ArgResolver
		name        string
		input       input.Service
		output      compiled.Service
		panic       string
	}{
		{
			name:   "db",
			input:  input.Service{Todo: true},
			output: compiled.Service{Name: "db", Todo: true},
		},
		{
			name: "db",
			input: input.Service{
				Constructor: `"mypkg/db".NewDB`,
				Fields: map[string]interface{}{
					"Port": "%port%",
					"Host": "%host%",
				},
				Tags: []input.Tag{
					{Name: "storage", Priority: 100},
					{Name: "sql", Priority: 0},
				},
			},
			output: compiled.Service{
				Name:        "db",
				Constructor: "alias.NewDB",
				Fields: []compiled.Field{
					{Name: "Host", Value: compiled.Arg{Code: "mock"}},
					{Name: "Port", Value: compiled.Arg{Code: "mock"}},
				},
				Tags: []compiled.Tag{
					{Name: "storage", Priority: 100},
					{Name: "sql"},
				},
			},
			argResolver: mockArgResolver{arg: compiled.Arg{Code: "mock"}},
		},
		{
			name: "db",
			input: input.Service{
				Args: []interface{}{
					"host",
				},
			},
			argResolver: mockArgResolver{error: fmt.Errorf("some error")},
			panic:       "service `db`: cannot resolve arg0: some error",
		},
		{
			name: "db",
			input: input.Service{
				Calls: []input.Call{
					{
						Method:    "SetHost",
						Args:      []interface{}{"%host%"},
						Immutable: false,
					},
				},
			},
			argResolver: mockArgResolver{error: fmt.Errorf("some error")},
			panic:       "service: `db`: call `SetHost`: cannot resolve arg0: some error",
		},
		{
			name: "db",
			input: input.Service{
				Fields: map[string]interface{}{
					"Host": "%host%",
				},
			},
			argResolver: mockArgResolver{error: fmt.Errorf("some error")},
			panic:       "service `db`: field `Host`: some error",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenaqrio #%d", i), func(t *testing.T) {
			if s.panic != "" {
				defer func() {
					r := recover()
					if !assert.NotEmpty(t, r) || !assert.Implements(t, (*error)(nil), r) {
						return
					}
					assert.EqualError(t, r.(error), s.panic)
				}()
			}
			compiler := Compiler{
				imports:     mockImports{alias: "alias"},
				argResolver: s.argResolver,
			}
			assert.Equal(t, s.output, compiler.handleService(s.name, s.input))
		})
	}
}

func TestCompiler_handleServiceType(t *testing.T) {
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
		Compiler{
			imports: mockImports{alias: "alias"},
		}.handleServiceType,
		scenarios...,
	)
}

func TestCompiler_handleServiceValue(t *testing.T) {
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
	}

	doTestInputOutput(
		t,
		Compiler{
			imports: mockImports{alias: "alias"},
		}.handleServiceValue,
		scenarios...,
	)
}

func TestCompiler_handleServiceConstructor(t *testing.T) {
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
		Compiler{
			imports: mockImports{alias: "alias"},
		}.handleServiceConstructor,
		scenarios...,
	)
}

type mockArgResolver struct {
	arg   compiled.Arg
	error error
}

func (m mockArgResolver) Resolve(interface{}) (compiled.Arg, error) {
	return m.arg, m.error
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
