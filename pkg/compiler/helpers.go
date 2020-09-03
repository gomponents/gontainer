package compiler

import (
	"github.com/gomponents/gontainer/pkg/syntax"
)

// deprecated
// use syntax.SanitizeImport
func sanitizeImport(i string) string {
	return syntax.SanitizeImport(i)
}
