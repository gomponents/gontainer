package compiler

import (
	"fmt"
	"regexp"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexMetaGoFn = regexp.MustCompile(`\A` + regex.MetaGoFn + `\z`)
)

type StepMeta struct {
	imports   ImportPrefixes
	tokenizer Tokenizer
}

func NewStepMeta(imports ImportPrefixes, tokenizer Tokenizer) *StepMeta {
	return &StepMeta{imports: imports, tokenizer: tokenizer}
}

func (s StepMeta) Do(i *input.DTO, result *compiled.DTO) error {
	result.Meta.Pkg = i.Meta.Pkg
	result.Meta.ContainerType = i.Meta.ContainerType

	for a, p := range i.Meta.Imports {
		err := s.imports.RegisterPrefix(a, sanitizeImport(p))
		if err != nil {
			return fmt.Errorf("cannot register alias: %s", err.Error())
		}
	}

	for fn, goFn := range i.Meta.Functions {
		_, m := regex.Match(regexMetaGoFn, goFn)
		s.tokenizer.RegisterFunction(sanitizeImport(m["import"]), m["fn"], fn)
	}

	return nil
}
