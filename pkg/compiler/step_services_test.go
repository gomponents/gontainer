package compiler

import (
	"fmt"
	"github.com/gomponents/gontainer/pkg/arguments"
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
			aliases: mockAliases{alias: "alias"},
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
			aliases: mockAliases{alias: "alias"},
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
			aliases: mockAliases{alias: "alias"},
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

	t.Run("Given error", func(t *testing.T) {
		argResolver := mockArgResolver{
			error: fmt.Errorf("resolver error"),
		}
		i := input.DTO{
			Services: map[string]input.Service{
				"db": {
					Constructor: "NewDB",
					Args:        []interface{}{"localhost"},
				},
			},
		}
		o := compiled.DTO{}
		err := StepServices{argResolver: argResolver}.Do(&i, &o)
		assert.EqualError(t, err, "service `db`: cannot resolve arg0: resolver error")
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

func TestStepServices_handleServiceFields(t *testing.T) {
	compiledFooParam := compiled.Arg{
		Code:            `container.GetParam("foo")`,
		Raw:             "%foo%",
		DependsOnParams: []string{"foo"},
	}

	scenarios := []struct {
		input    map[string]interface{}
		output   []compiled.Field
		error    string
		resolver arguments.Resolver
	}{
		{
			input: map[string]interface{}{
				"Host": "localhost",
			},
			error:    "field `Host`: my resolver error",
			resolver: mockArgResolver{error: fmt.Errorf("my resolver error")},
		},
		{
			input: map[string]interface{}{
				"Port":  3306,
				"Host":  "localhot",
				"Debug": true,
			},
			resolver: mockArgResolver{arg: compiledFooParam},
			output: []compiled.Field{
				{Name: "Debug", Value: compiledFooParam},
				{Name: "Host", Value: compiledFooParam},
				{Name: "Port", Value: compiledFooParam},
			},
		},
		{
			input:  nil,
			output: nil,
			error:  "",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			o, err := StepServices{argResolver: s.resolver}.handleServiceFields(s.input)
			if s.error != "" {
				assert.EqualError(t, err, s.error)
				assert.Nil(t, o)
				return
			}
			assert.Equal(t, s.output, o)
		})
	}
}

func TestStepServices_handleServiceCalls(t *testing.T) {
	compiledFooParam := compiled.Arg{
		Code:            `container.GetParam("foo")`,
		Raw:             "%foo%",
		DependsOnParams: []string{"foo"},
	}

	scenarios := []struct {
		input    []input.Call
		output   []compiled.Call
		error    string
		resolver arguments.Resolver
	}{
		{
			input:  nil,
			output: nil,
			error:  "",
		},
		{
			input: []input.Call{
				{Method: "SetFoo", Args: []interface{}{"foo"}, Immutable: true},
				{Method: "SetBar", Args: []interface{}{"bar"}, Immutable: false},
			},
			output: []compiled.Call{
				{Method: "SetFoo", Args: []compiled.Arg{compiledFooParam}, Immutable: true},
				{Method: "SetBar", Args: []compiled.Arg{compiledFooParam}, Immutable: false},
			},
			resolver: mockArgResolver{arg: compiledFooParam},
		},
		{
			input: []input.Call{
				{Method: "SetFoo", Args: []interface{}{"foo"}, Immutable: true},
				{Method: "SetBar", Args: []interface{}{"bar"}, Immutable: false},
			},
			error:    "call `SetFoo`: cannot resolve arg0: fancy error",
			resolver: mockArgResolver{error: fmt.Errorf("fancy error")},
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			o, err := StepServices{argResolver: s.resolver}.handleServiceCalls(s.input)
			if s.error != "" {
				assert.EqualError(t, err, s.error)
				assert.Nil(t, o)
				return
			}
			assert.Equal(t, s.output, o)
		})
	}
}

type mockAliases struct {
	alias string
}

func (m mockAliases) GetAlias(string) string {
	return m.alias
}

type mockArgResolver struct {
	arg   compiled.Arg
	error error
}

func (m mockArgResolver) Resolve(interface{}) (compiled.Arg, error) {
	return m.arg, m.error
}
