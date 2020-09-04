package cmd

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

func newDefaultCompiler(imports imports.Imports) *compiler.Compiler {
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
		compiler.NewStepValidateInput(input.NewDefaultValidator()),                     // validate input
		compiler.NewStepMeta(imports, tokenizer),                                       // process meta
		compiler.NewStepParams(paramResolver),                                          // process params
		compiler.NewStepServices(imports, arguments.NewDefaultResolver(paramResolver)), // process services
		compiler.NewStepValidateOutput(compiled.NewDefaultValidator()),                 // validate output
	)
}
