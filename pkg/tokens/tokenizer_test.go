package tokens

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatternTokenizer_Tokenize(t *testing.T) {
	tokenizer := NewPatternTokenizer(
		[]TokenFactoryStrategy{
			TokenPercentSign{},
			TokenReference{},
			TokenString{},
		},
		nil,
	)

	scenarios := []struct {
		input  string
		output []Token
		error  string
	}{
		{
			input: "%name%",
			output: []Token{
				{
					Kind:      KindReference,
					Raw:       "%name%",
					DependsOn: []string{"name"},
					Code:      "",
				},
			},
		},
		{
			input: "%firstname% %lastname%",
			output: []Token{
				{
					Kind:      KindReference,
					Raw:       "%firstname%",
					DependsOn: []string{"firstname"},
					Code:      "",
				},
				{
					Kind: KindString,
					Raw:  " ",
				},
				{
					Kind:      KindReference,
					Raw:       "%lastname%",
					DependsOn: []string{"lastname"},
					Code:      "",
				},
			},
		},
		{
			input: "%%",
			output: []Token{
				{
					Kind: KindCode,
					Raw:  "%%",
					Code: `"%"`,
				},
			},
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			tkns, err := tokenizer.Tokenize(s.input)
			if s.error == "" {
				assert.NoError(t, err)
				assert.Equal(t, s.output, tkns)
				return
			}

			assert.Empty(t, tkns)
			assert.EqualError(t, err, s.error)
		})
	}
}
