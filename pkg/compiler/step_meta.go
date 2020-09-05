package compiler

import (
	"fmt"
	"regexp"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/regex"
	"github.com/gomponents/gontainer/pkg/tokens"
)

var (
	regexMetaGoFn = regexp.MustCompile(`\A` + regex.MetaGoFn + `\z`)
)

type StepMeta struct {
	imports   imports.Prefixes
	functions tokens.Functions
}

func NewStepMeta(imports imports.Prefixes, functions tokens.Functions) *StepMeta {
	return &StepMeta{imports: imports, functions: functions}
}

func (s StepMeta) Do(i input.DTO, result *compiled.DTO) error {
	result.Meta.Pkg = i.Meta.Pkg
	result.Meta.ContainerType = i.Meta.ContainerType

	if err := s.handleImports(i, result); err != nil {
		return err
	}

	s.handleFunctions(i)

	return nil
}

func (s StepMeta) handleImports(i input.DTO, result *compiled.DTO) error {
	for a, p := range i.Meta.Imports {
		err := s.imports.RegisterPrefix(a, sanitizeImport(p))
		if err != nil {
			return fmt.Errorf("cannot register alias: %s", err.Error())
		}
	}
	return nil
}

func (s StepMeta) handleFunctions(i input.DTO) {
	for fn, goFn := range i.Meta.Functions {
		_, m := regex.Match(regexMetaGoFn, goFn)
		s.functions.RegisterFunction(sanitizeImport(m["import"]), m["fn"], fn)
	}
}
