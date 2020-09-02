package defaults

import (
	"testing"

	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultCompiler(t *testing.T) {
	assert.NotEmpty(t, NewDefaultCompiler(mockImports{}))
}

type mockImports struct {
}

func (m mockImports) GetAlias(string) string {
	return ""
}

func (m mockImports) GetImports() []imports.Import {
	return nil
}

func (m mockImports) RegisterPrefix(shortcut string, path string) error {
	return nil
}
