package tokens

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_toExpr(t *testing.T) {
	scenarios := []struct {
		input  string
		output string
		ok     bool
	}{
		{"%", "", false},
		{"%%", "", true},
		{"%hello%", "hello", true},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			output, ok := toExpr(s.input)
			assert.Equal(t, s.output, output)
			assert.Equal(t, s.ok, ok)
		})
	}
}

func TestTokenPercentSign_Supports(t *testing.T) {
	t.Run("Given supported scenario", func(t *testing.T) {
		assert.True(t, TokenPercentSign{}.Supports("%%"))
	})
	t.Run("Given unsupported scenario", func(t *testing.T) {
		assert.False(t, TokenPercentSign{}.Supports("hello"))
	})
}

func TestTokenPercentSign_Create(t *testing.T) {
	tkn, err := TokenPercentSign{}.Create("%%")
	assert.NoError(t, err)
	assert.Equal(
		t,
		Token{
			Kind: KindCode,
			Raw:  "%%",
			Code: `"%"`,
		},
		tkn,
	)
}
