package template

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer-helpers/caller"
	"github.com/stretchr/testify/assert"
)

func Test_createDefaultFunctions(t *testing.T) {
	fncs := createDefaultFunctions(mockAliases{"myAlias"})

	scenarios := []struct {
		fn     string
		input  []interface{}
		output string
	}{
		{
			fn:     "export",
			input:  []interface{}{5},
			output: "5",
		},
		{
			fn:     "importAlias",
			input:  []interface{}{"my/package.name"},
			output: "myAlias",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			o := caller.MustCall(fncs[s.fn], s.input...)
			assert.Equal(t, s.output, o[0])
		})
	}
}

type mockAliases struct {
	alias string
}

func (m mockAliases) GetAlias(string) string {
	return m.alias
}
