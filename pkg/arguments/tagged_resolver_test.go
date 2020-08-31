package arguments

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/stretchr/testify/assert"
)

func TestTaggedResolver_Supports(t *testing.T) {
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
			input:  "!tagged tagName",
			output: true,
		},
		{
			input:  "!tagged tag_name",
			output: true,
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			assert.Equal(
				t,
				s.output,
				NewTaggedResolver().Supports(s.input),
			)
		})
	}
}

func TestTaggedResolver_Resolve(t *testing.T) {
	arg, err := NewTaggedResolver().Resolve("!tagged my_tag")
	assert.NoError(t, err)
	assert.Equal(
		t,
		compiled.Arg{
			Code:          `container.MustGetByTag("my_tag")`,
			Raw:           "!tagged my_tag",
			DependsOnTags: []string{"my_tag"},
		},
		arg,
	)
}
