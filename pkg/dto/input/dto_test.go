package input

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestCall_UnmarshalYAML(t *testing.T) {
	// todo more tests
	y := `
- ["SetName", ["Mary"], true]
`
	calls := make([]Call, 0)
	err := yaml.Unmarshal([]byte(y), &calls)
	assert.NoError(t, err)
	assert.Equal(
		t,
		[]Call{
			{
				Method:    "SetName",
				Args:      []interface{}{"Mary"},
				Immutable: true,
			},
		},
		calls,
	)
}

func TestCreateDefaultDTO(t *testing.T) {
	assert.NoError(
		t,
		NewDefaultValidator().Validate(CreateDefaultDTO()),
	)
}

func TestTag_UnmarshalYAML(t *testing.T) {
	scenarios := []struct {
		input  string
		output Tag
		error  string
	}{
		{
			input: `
name: hello
priority: 100
`,
			output: Tag{Name: "hello", Priority: 100},
		},
		{
			input: `
name: hello
`,
			output: Tag{Name: "hello"},
		},
		{
			input: `
name: 123
`,
			error: "name must be an instance of string",
		},
		{
			input: `
name: hello
priority: high
`,
			error: "priority must be an instance of int",
		},
		{
			input: `
priority: 100
`,
			error: "name of tag is required",
		},
		{
			input:  "tag",
			output: Tag{Name: "tag"},
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			tag := Tag{}
			err := yaml.Unmarshal([]byte(s.input), &tag)
			if s.error != "" {
				assert.EqualError(t, err, s.error)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, s.output, tag)
		})
	}
}
