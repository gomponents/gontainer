package compiler

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/stretchr/testify/assert"
)

func TestCompiler_Compile(t *testing.T) {
	t.Run("Given error", func(t *testing.T) {
		c := NewCompiler(stepMock{error: fmt.Errorf("test error")})
		_, err := c.Compile(input.DTO{})
		assert.EqualError(t, err, "test error")
	})
	t.Run("No errors", func(t *testing.T) {
		c := NewCompiler(stepMock{error: nil})
		_, err := c.Compile(input.DTO{})
		assert.NoError(t, err)
	})
}

type stepMock struct {
	error error
}

func (s stepMock) Do(input.DTO, *compiled.DTO) error {
	return s.error
}
