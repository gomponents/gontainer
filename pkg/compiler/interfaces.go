package compiler

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
)

type InputValidator interface {
	Validate(input.DTO) error
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
