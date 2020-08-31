package input

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestCall_UnmarshalYAML(t *testing.T) {
	scenarios := []struct {
		input  string
		output Call
		error  string
	}{
		{
			input: `
- SetName
- ["Mary"]
- true`,
			output: Call{
				Method:    "SetName",
				Args:      []interface{}{"Mary"},
				Immutable: true,
			},
		},
		{
			input: `[SetName, ["Mary"], false]`,
			output: Call{
				Method:    "SetName",
				Args:      []interface{}{"Mary"},
				Immutable: false,
			},
		},
		{
			input: `
- SetName
- ["Mary"]
- 115`,
			error: "third element of object Call must be a bool, `int` given",
		},
		{
			input: `
- SetName
- false
- true`,
			error: "second element of object Call must be an array, `bool` given",
		},
		{
			input: `[78, ["Mary"], false]`,
			error: "first element of object Call must be a string, `int` given",
		},
		{
			input: `[1, 2, 3, 4]`,
			error: "object Call must contain 1 - 3 args, 4 given",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			call := Call{}
			err := yaml.Unmarshal([]byte(s.input), &call)
			if s.error != "" {
				assert.EqualError(t, err, s.error)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, s.output, call)
		})
	}
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
