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

func TestTokenFactoryStrategy(t *testing.T) {
	scenarios := []struct {
		factory  TokenFactoryStrategy
		expr     string
		supports bool
		token    Token
	}{
		{
			factory:  TokenPercentSign{},
			expr:     "%%",
			supports: true,
			token: Token{
				Kind: KindCode,
				Raw:  "%%",
				Code: `"%"`,
			},
		},
		{
			factory:  TokenPercentSign{},
			expr:     "%",
			supports: false,
		},
		{
			factory:  TokenReference{},
			expr:     "%name%",
			supports: true,
			token: Token{
				Kind:      KindReference,
				Raw:       "%name%",
				DependsOn: []string{"name"},
			},
		},
		{
			factory:  TokenReference{},
			expr:     "%name",
			supports: false,
		},
		{
			factory:  NewTokenSimpleFunction(mockAliases{alias: "osAlias"}, "env", "os", "Getenv"),
			expr:     `%env("FOO")%`,
			supports: true,
			token: Token{
				Kind: KindCode,
				Raw:  `%env("FOO")%`,
				Code: `osAlias.WrapMustCallProvider("cannot execute %env(\"FOO\")%", osAlias.Getenv, "FOO")`,
			},
		},
		{
			factory:  NewTokenSimpleFunction(mockAliases{alias: "osAlias"}, "env", "os", "Getenv"),
			expr:     `lorep ipsum`,
			supports: false,
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			if s.supports {
				if !assert.True(t, s.factory.Supports(s.expr)) {
					return
				}

				token, err := s.factory.Create(s.expr)

				assert.NoError(t, err)
				assert.Equal(t, s.token, token)

				return
			}

			assert.False(t, s.factory.Supports(s.expr))
		})
	}
}
