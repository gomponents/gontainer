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
		panic  bool
	}{
		{
			fn:     "export",
			input:  []interface{}{5},
			output: "5",
		},
		{
			fn:    "export",
			input: []interface{}{struct{}{}},
			panic: true, // parameter of type `struct {}` is not supported
		},
		{
			fn:     "importAlias",
			input:  []interface{}{"my/package.name"},
			output: "myAlias",
		},
		{
			fn:     "callerAlias",
			output: "myAlias",
		},
		{
			fn:     "containerAlias",
			output: "myAlias",
		},
		{
			fn:     "setterAlias",
			output: "myAlias",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			defer func() {
				if s.panic {
					assert.NotNil(t, recover())
					return
				}

				assert.Nil(t, recover())
			}()

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
