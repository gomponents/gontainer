package input

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateParams(t *testing.T) {
	scenarios := []struct {
		name  string
		val   interface{}
		error string
	}{
		{
			name: "test",
			val:  5,
		},
		{
			name:  "@",
			error: "invalid parameter name `@`",
		},
		{
			name:  "param",
			val:   struct{}{},
			error: "unsupported type `struct {}` of parameter `param`",
		},
	}
	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			d := DTO{}
			d.Params = map[string]interface{}{
				s.name: s.val,
			}
			err := ValidateParams(d)

			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, s.error)
		})
	}
}
