package pkg

import (
	"github.com/gomponents/gontainer-helpers/exporters"
	"github.com/gomponents/gontainer/pkg/arguments"
	"github.com/gomponents/gontainer/pkg/compiler"
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/parameters"
	"github.com/gomponents/gontainer/pkg/tokens"
)

func NewDefaultCompiler(imports imports.Imports) *compiler.Compiler {
	tokenizer := tokens.NewPatternTokenizer(
		[]tokens.TokenFactoryStrategy{
			tokens.TokenPercentSign{},
			tokens.TokenReference{},
			tokens.TokenString{},
		},
		imports,
	)

	paramResolver := parameters.NewSimpleResolver(
		tokenizer,
		exporters.NewDefaultExporter(),
		imports,
	)

	return compiler.NewCompiler(
		input.NewDefaultValidator(),
		compiled.NewDefaultValidator(),
		imports,
		tokenizer,
		paramResolver,
		arguments.NewDefaultResolver(paramResolver),
	)
}
