package cmd

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
)

type cmplr interface {
	Compile(i input.DTO) (compiled.DTO, error)
}

type compile struct {
	compiler cmplr
}

func newCompile(compiler cmplr) *compile {
	return &compile{compiler: compiler}
}

func (c compile) run(i *input.DTO, o *compiled.DTO) (err error) {
	*o, err = c.compiler.Compile(*i)
	return
}
