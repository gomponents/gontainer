package tokens

import (
	"fmt"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/regex"
)

type Kind uint

const (
	KindString Kind = iota
	KindReference
	KindCode
)

const (
	TokenDelimiter = "%"
)

var (
	regexTokenRef = regex.MustCompileAz(`[a-zA-Z][a-zA-Z0-9_]*((\.)[a-zA-Z0-9_]+)*`)
	regexSimpleFn = regex.MustCompileAz(`(?P<fn>[a-zA-Z][a-zA-Z0-9_]*((\.)[a-zA-Z0-9_]+)*)\((?P<params>.*)\)`)
)

type Token struct {
	Kind      Kind
	Raw       string
	DependsOn []string
	Code      string
}

type TokenFactoryStrategy interface {
	Supports(expr string) bool
	Create(expr string) (Token, error)
}

// toExpr removes surrounding delimiters
func toExpr(expr string) (string, bool) {
	runes := []rune(expr)
	if len(runes) < 2 {
		return "", false
	}

	if string(runes[0]) != TokenDelimiter || string(runes[len(runes)-1]) != TokenDelimiter {
		return "", false
	}

	return string(runes[1 : len(runes)-1]), true
}

// %%
type TokenPercentSign struct{}

func (t TokenPercentSign) Supports(expr string) bool {
	return expr == "%%"
}

func (t TokenPercentSign) Create(expr string) (Token, error) {
	return Token{
		Kind:      KindCode,
		Raw:       "%%",
		DependsOn: nil,
		Code:      `"%"`,
	}, nil
}

// %my.param%
type TokenReference struct{}

func (t TokenReference) Supports(s string) bool {
	expr, ok := toExpr(s)

	return ok && regexTokenRef.MatchString(expr)
}

func (t TokenReference) Create(s string) (Token, error) {
	ref, _ := toExpr(s)

	return Token{
		Kind:      KindReference,
		Raw:       s,
		DependsOn: []string{ref},
	}, nil
}

type TokenString struct{}

func (t TokenString) Supports(expr string) bool {
	return true
}

func (t TokenString) Create(expr string) (Token, error) {
	return Token{
		Kind:      KindString,
		Raw:       expr,
		DependsOn: nil,
	}, nil
}

// %env(ENV_VAR)%
type TokenSimpleFunction struct {
	aliases  imports.Aliases
	fn       string
	goImport string
	goFn     string
}

func NewTokenSimpleFunction(aliases imports.Aliases, fn string, goImport string, goFn string) *TokenSimpleFunction {
	return &TokenSimpleFunction{aliases: aliases, fn: fn, goImport: goImport, goFn: goFn}
}

func (t TokenSimpleFunction) Supports(expr string) bool {
	e, ok := toExpr(expr)
	if !ok {
		return false
	}

	ok, m := regex.Match(regexSimpleFn, e)
	return ok && m["fn"] == t.fn
}

func (t TokenSimpleFunction) Create(expr string) (Token, error) {
	e, _ := toExpr(expr)
	_, m := regex.Match(regexSimpleFn, e)
	fn := fmt.Sprintf("%s(%s)", t.goFn, m["params"])
	if t.goImport != "" && t.goImport != `"."` {
		fn = fmt.Sprintf("%s.%s", t.aliases.GetAlias(t.goImport), fn)
	}
	return Token{
		Kind: KindCode,
		Raw:  expr,
		Code: fn,
	}, nil
}
