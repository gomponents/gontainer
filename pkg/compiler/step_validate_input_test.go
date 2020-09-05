package compiler

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/stretchr/testify/assert"
)

func TestStepValidateInput_Do(t *testing.T) {
	t.Run("Given error", func(t *testing.T) {
		s := NewStepValidateInput(mockInputValidator{error: fmt.Errorf("input error")})
		assert.EqualError(t, s.Do(input.DTO{}, &compiled.DTO{}), "input error")
	})
	t.Run("No errors", func(t *testing.T) {
		s := NewStepValidateInput(mockInputValidator{error: nil})
		assert.NoError(t, s.Do(input.DTO{}, &compiled.DTO{}))
	})
}

type mockInputValidator struct {
	error error
}

func (m mockInputValidator) Validate(input.DTO) error {
	return m.error
}
