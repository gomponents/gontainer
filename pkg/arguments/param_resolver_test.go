package arguments

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/parameters"
	"github.com/stretchr/testify/assert"
)

func TestParamResolver_Supports(t *testing.T) {
	assert.True(t, ParamResolver{}.Supports("whatever"))
}

func TestParamResolver_Resolve(t *testing.T) {
	a := complex128(0.588 + 0i)
	fmt.Println(a)
	expectedErr := fmt.Errorf("test error")
	m := mockResolver{
		expr: parameters.Expr{
			Code:      `host + port`,
			DependsOn: []string{"host", "port"},
		},
		err: expectedErr,
	}

	arg, err := NewParamResolver(m).Resolve("%host%:%port%")
	assert.Same(t, expectedErr, err)
	assert.Equal(
		t,
		compiled.Arg{
			Code:              "host + port",
			Raw:               "%host%:%port%",
			DependsOnParams:   []string{"host", "port"},
			DependsOnServices: nil,
		},
		arg,
	)
}

type mockResolver struct {
	expr parameters.Expr
	err  error
}

func (m mockResolver) Resolve(v interface{}) (parameters.Expr, error) {
	r := m.expr
	r.Raw = v
	return r, m.err
}
