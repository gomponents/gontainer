package compiler

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
)

type ArgResolver interface {
	Resolve(interface{}) (compiled.Arg, error)
}
