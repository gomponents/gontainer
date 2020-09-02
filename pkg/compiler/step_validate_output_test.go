package compiler

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/stretchr/testify/assert"
)

func TestStepValidateOutput_Do(t *testing.T) {
	t.Run("Given error", func(t *testing.T) {
		s := NewStepValidateOutput(mockCompiledValidator{error: fmt.Errorf("compiler error")})
		assert.EqualError(t, s.Do(&input.DTO{}, &compiled.DTO{}), "compiler error")
	})
	t.Run("No errors", func(t *testing.T) {
		s := NewStepValidateOutput(mockCompiledValidator{error: nil})
		assert.NoError(t, s.Do(&input.DTO{}, &compiled.DTO{}))
	})
}

type mockCompiledValidator struct {
	error error
}

func (m mockCompiledValidator) Validate(compiled.DTO) error {
	return m.error
}
