package tokens

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatternTokenizer_Tokenize(t *testing.T) {
	defaultTokenizer := NewPatternTokenizer(
		[]TokenFactoryStrategy{
			TokenPercentSign{},
			TokenReference{},
			TokenString{},
		},
		nil,
	)
	emptyTokenizer := NewPatternTokenizer(nil, nil)

	scenarios := []struct {
		tokenizer Tokenizer
		input     string
		output    []Token
		error     string
	}{
		{
			tokenizer: defaultTokenizer,
			input:     "%name%",
			output: []Token{
				{
					Kind:      KindReference,
					Raw:       "%name%",
					DependsOn: []string{"name"},
				},
			},
		},
		{
			tokenizer: defaultTokenizer,
			input:     "%firstname% %lastname%",
			output: []Token{
				{
					Kind:      KindReference,
					Raw:       "%firstname%",
					DependsOn: []string{"firstname"},
				},
				{
					Kind: KindString,
					Raw:  " ",
				},
				{
					Kind:      KindReference,
					Raw:       "%lastname%",
					DependsOn: []string{"lastname"},
				},
			},
		},
		{
			tokenizer: defaultTokenizer,
			input:     "%%",
			output: []Token{
				{
					Kind: KindCode,
					Raw:  "%%",
					Code: `"%"`,
				},
			},
		},
		{
			tokenizer: defaultTokenizer,
			input:     "John %",
			output: []Token{
				{
					Kind: KindString,
					Raw:  "John ",
				},
				{
					Kind: KindString,
					Raw:  "%",
				},
			},
		},
		{
			tokenizer: emptyTokenizer,
			input:     "test",
			error:     "invalid token `test`",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			tkns, err := s.tokenizer.Tokenize(s.input)
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
