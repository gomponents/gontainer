package arguments

import (
	"fmt"
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceResolver_Supports(t *testing.T) {
	scenarios := []struct {
		input  interface{}
		output bool
	}{
		{
			input:  5,
			output: false,
		},
		{
			input:  "service",
			output: false,
		},
		{
			input:  "@service",
			output: true,
		},
		{
			input:  "@@",
			output: false,
		},
		{
			input:  "@",
			output: false,
		},
		{
			input:  "@ service",
			output: false,
		},
		{
			input:  "@s.e_r.v_i_c_e",
			output: true,
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			assert.Equal(
				t,
				s.output,
				NewServiceResolver().Supports(s.input),
			)
		})
	}
}

func TestServiceResolver_Resolve(t *testing.T) {
	arg, err := ServiceResolver{}.Resolve("@db")
	assert.NoError(t, err)
	assert.Equal(
		t,
		compiled.Arg{
			Code:              `container.MustGet("db")`,
			Raw:               "@db",
			DependsOnServices: []string{"db"},
		},
		arg,
	)
}
