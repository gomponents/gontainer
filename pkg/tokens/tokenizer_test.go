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
	errorTokenizer := NewPatternTokenizer(
		[]TokenFactoryStrategy{
			tokenMock{
				supports: true,
				token:    Token{},
				error:    fmt.Errorf("custom error"),
			},
		}, nil,
	)

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
		{
			tokenizer: errorTokenizer,
			input:     "hello",
			error:     "cannot create token `hello`: custom error",
		},
		{
			tokenizer: errorTokenizer,
			input:     "name: %name%",
			error:     "cannot create token `name: `: custom error",
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

func TestPatternTokenizer_RegisterFunction(t *testing.T) {
	tokenizer := NewPatternTokenizer(nil, mockAliasProvider{alias: "myAlias"})
	assert.Empty(t, tokenizer.strategies)
	const expr = `%env("VAR")%`

	tkns1, err1 := tokenizer.Tokenize(expr)
	assert.EqualError(t, err1, "invalid token `%env(\"VAR\")%`")
	assert.Empty(t, tkns1)

	tokenizer.RegisterFunction("os", "Getenv", "env")
	assert.NotEmpty(t, tokenizer.strategies)
	tkns2, err2 := tokenizer.Tokenize(expr)
	assert.NoError(t, err2)
	assert.Equal(
		t,
		[]Token{
			{
				Kind: KindCode,
				Raw:  expr,
				Code: `myAlias.Getenv("VAR")`,
			},
		},
		tkns2,
	)
}

type tokenMock struct {
	supports bool
	token    Token
	error    error
}

func (t tokenMock) Supports(string) bool {
	return t.supports
}

func (t tokenMock) Create(string) (Token, error) {
	return t.token, t.error
}
