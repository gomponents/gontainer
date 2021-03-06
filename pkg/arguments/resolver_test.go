package arguments

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultResolver(t *testing.T) {
	resolver := NewDefaultResolver(nil, nil)
	assert.NotEmpty(t, resolver.subResolvers)
}

func TestResolver_Resolve(t *testing.T) {
	scenarios := []struct {
		resolver *ChainResolver
		input    interface{}
		output   compiled.Arg
		error    string
	}{
		{
			resolver: NewChainResolver(NewServiceResolver()),
			input:    "@db",
			output: compiled.Arg{
				Code:              `container.MustGet("db")`,
				Raw:               "@db",
				DependsOnServices: []string{"db"},
			},
		},
		{
			resolver: NewChainResolver(),
			input:    "%name%",
			error:    "cannot resolve argument `%name%`",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			arg, err := s.resolver.Resolve(s.input)
			if s.error == "" {
				assert.NoError(t, err)
				assert.Equal(t, s.output, arg)
				return
			}
			assert.EqualError(t, err, s.error)
		})
	}
}
