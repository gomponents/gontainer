package input

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isPrimitiveType(t *testing.T) {
	scenarios := []struct {
		input  interface{}
		output bool
	}{
		{nil, true},
		{5, true},
		{math.Pi, true},
		{false, true},
		{struct{}{}, false},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			assert.Equal(t, s.output, isPrimitiveType(s.input))
		})
	}
}
