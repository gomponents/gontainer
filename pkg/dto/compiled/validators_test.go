package compiled

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultValidator(t *testing.T) {
	assert.NotEmpty(t, NewDefaultValidator().validators)
}

func TestChainValidator_Validate(t *testing.T) {
	t.Run("Given error", func(t *testing.T) {
		v := NewChainValidator(func(DTO) error {
			return fmt.Errorf("my custom error")
		})
		assert.EqualError(
			t,
			v.Validate(DTO{}),
			"my custom error",
		)
	})
	t.Run("No error", func(t *testing.T) {
		v := NewChainValidator(func(DTO) error {
			return nil
		})
		assert.NoError(t, v.Validate(DTO{}))
	})
}
