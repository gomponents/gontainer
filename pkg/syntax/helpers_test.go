package syntax

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSanitizeImport(t *testing.T) {
	scenarios := []struct {
		input  string
		output string
	}{
		{
			input:  `"pkg"`,
			output: "pkg",
		},
		{
			input:  "config",
			output: "config",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			assert.Equal(t, s.output, SanitizeImport(s.input))
		})
	}
}
