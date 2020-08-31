package tokens

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_toExpr(t *testing.T) {
	scenarios := []struct {
		input  string
		output string
		ok     bool
	}{
		{"%", "", false},
		{"%%", "", true},
		{"%hello%", "hello", true},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			output, ok := toExpr(s.input)
			assert.Equal(t, s.output, output)
			assert.Equal(t, s.ok, ok)
		})
	}
}
