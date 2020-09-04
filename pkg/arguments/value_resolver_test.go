package arguments

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/stretchr/testify/assert"
)

func TestValueResolver_Supports(t *testing.T) {
	scenarios := []struct {
		input  interface{}
		output bool
	}{
		{
			input:  5,
			output: false,
		},
		{
			input:  "GlobalVariable",
			output: false,
		},
		{
			input:  "!value GlobalVariable",
			output: true,
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			assert.Equal(
				t,
				s.output,
				ValueResolver{}.Supports(s.input),
			)
		})
	}
}

func TestValueResolver_Resolve(t *testing.T) {
	scenarios := []struct {
		input   string
		output  compiled.Arg
		aliases imports.Aliases
	}{
		{
			input:  "!value &Logger{}",
			output: compiled.Arg{Code: "&Logger{}", Raw: "!value &Logger{}"},
		},
		{
			input:  `!value ".".GlobalConfig.DbConfig`,
			output: compiled.Arg{Code: "GlobalConfig.DbConfig", Raw: `!value ".".GlobalConfig.DbConfig`},
		},
		{
			input:   `!value "config".GlobalConfig.DbConfig`,
			output:  compiled.Arg{Code: "alias123.GlobalConfig.DbConfig", Raw: `!value "config".GlobalConfig.DbConfig`},
			aliases: mockAliases{alias: "alias123"},
		},
		{
			input:   `!value "config".GlobalConfig.DbConfig`,
			output:  compiled.Arg{Code: "alias123.GlobalConfig.DbConfig", Raw: `!value "config".GlobalConfig.DbConfig`},
			aliases: mockAliases{alias: "alias123"},
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenaqqrio #%d", i), func(t *testing.T) {
			arg, err := NewValueResolver(s.aliases).Resolve(s.input)
			assert.NoError(t, err)
			assert.Equal(t, s.output, arg)
		})
	}
}

type mockAliases struct {
	alias string
}

func (m mockAliases) GetAlias(string) string {
	return m.alias
}
