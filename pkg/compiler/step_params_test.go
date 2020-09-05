package compiler

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/gomponents/gontainer/pkg/parameters"
	"github.com/stretchr/testify/assert"
)

func TestStepParams_Do(t *testing.T) {
	fooParam := parameters.Expr{
		Code:      `container.Get("foo")`,
		Raw:       "$foo%",
		DependsOn: []string{"foo"},
	}

	scenarios := []struct {
		input    map[string]interface{}
		output   []compiled.Param
		error    string
		resolver parameters.Resolver
	}{
		{
			input: map[string]interface{}{
				"Port": 3306,
				"Host": "localhost",
			},
			error:    "cannot resolve param `Host`: cannot resolve given param",
			resolver: mockParamResolver{error: fmt.Errorf("cannot resolve given param")},
		},
		{
			input: map[string]interface{}{
				"Port":  3306,
				"Debug": false,
				"Host":  "localhost",
			},
			output: []compiled.Param{
				{Name: "Debug", Code: fooParam.Code, Raw: fooParam.Raw, DependsOn: fooParam.DependsOn},
				{Name: "Host", Code: fooParam.Code, Raw: fooParam.Raw, DependsOn: fooParam.DependsOn},
				{Name: "Port", Code: fooParam.Code, Raw: fooParam.Raw, DependsOn: fooParam.DependsOn},
			},
			resolver: mockParamResolver{expr: fooParam},
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			c := compiled.DTO{}
			err := NewStepParams(s.resolver).Do(input.DTO{Params: s.input}, &c)
			if s.error != "" {
				assert.EqualError(t, err, s.error)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, s.output, c.Params)
		})
	}
}

type mockParamResolver struct {
	expr  parameters.Expr
	error error
}

func (m mockParamResolver) Resolve(interface{}) (parameters.Expr, error) {
	return m.expr, m.error
}
