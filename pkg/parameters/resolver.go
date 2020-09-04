package parameters

import (
	"fmt"
	"strings"

	"github.com/gomponents/gontainer-helpers/exporters"
	"github.com/gomponents/gontainer/pkg/consts"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/tokens"
)

type Resolver interface {
	Resolve(interface{}) (Expr, error)
}

type SimpleResolver struct {
	tokenizer tokens.Tokenizer
	exporter  exporters.Exporter
	aliases   imports.Aliases
}

func NewSimpleResolver(tokenizer tokens.Tokenizer, exporter exporters.Exporter, aliases imports.Aliases) *SimpleResolver {
	return &SimpleResolver{tokenizer: tokenizer, exporter: exporter, aliases: aliases}
}

type Expr struct {
	Code      string
	Raw       interface{}
	DependsOn []string
}

func (s SimpleResolver) Resolve(v interface{}) (Expr, error) {
	if str, ok := v.(string); ok {
		return s.resolveString(str)
	}

	return s.resolveNonString(v)
}

func (s SimpleResolver) resolveNonString(v interface{}) (Expr, error) {
	code, err := s.exporter.Export(v)
	if err != nil {
		return Expr{}, err
	}

	return Expr{
		Code:      code,
		Raw:       v,
		DependsOn: nil,
	}, nil
}

func (s SimpleResolver) resolveString(v string) (Expr, error) {
	tkns, tokenizerErr := s.tokenizer.Tokenize(v)
	if tokenizerErr != nil {
		return Expr{}, tokenizerErr
	}

	var dependsOn []string

	solveTokenCode := func(t tokens.Token) (string, error) {
		switch t.Kind {
		case tokens.KindString:
			return s.exporter.Export(t.Raw)
		case tokens.KindReference:
			runes := []rune(t.Raw)
			depID := string(runes[1 : len(runes)-1])
			dependsOn = append(dependsOn, depID)
			return fmt.Sprintf("container.MustGetParam(%+q)", depID), nil
		case tokens.KindCode:
			return t.Code, nil
		default:
			return "", fmt.Errorf("unexpected Kind %v", t.Kind)
		}
	}

	if len(tkns) == 0 {
		return Expr{}, fmt.Errorf("unexpected error, tokenizer returned 0 tokens")
	}

	if len(tkns) == 1 {
		t := tkns[0]
		code, solveErr := solveTokenCode(t)
		if solveErr != nil {
			return Expr{}, solveErr
		}

		return Expr{
			Code:      code,
			Raw:       v,
			DependsOn: dependsOn,
		}, nil
	}

	codeParts := make([]string, 0)
	for _, t := range tkns {
		tmp, solveErr := solveTokenCode(t)
		if solveErr != nil {
			return Expr{}, solveErr
		}

		tmp = s.aliases.GetAlias(consts.GontainerHelperPath+"/exporters") + `.MustToString(` + tmp + `)`
		codeParts = append(codeParts, tmp)
	}

	return Expr{
		Code:      strings.Join(codeParts, " + "),
		Raw:       v,
		DependsOn: dependsOn,
	}, nil
}
