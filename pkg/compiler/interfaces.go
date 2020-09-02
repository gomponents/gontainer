package compiler

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
)

type InputValidator interface {
	Validate(input.DTO) error
}

// todo rename it
type ImportAliases interface {
	GetAlias(string) string
}

// todo rename it
type ImportPrefixes interface {
	RegisterPrefix(shortcut string, path string) error
}

type Functions interface {
	RegisterFunction(goImport string, goFunc string, tokenFunc string)
}

type CompiledValidator interface {
	Validate(compiled.DTO) error
}

type ArgResolver interface {
	Resolve(interface{}) (compiled.Arg, error)
}
