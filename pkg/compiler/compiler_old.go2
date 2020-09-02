package compiler

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/gomponents/gontainer/pkg/parameters"
	"github.com/gomponents/gontainer/pkg/regex"
)

type Imports interface {
	GetAlias(string) string
	RegisterPrefix(shortcut string, path string) error
}

type CompilerOld struct {
	inputValidator    InputValidator
	compiledValidator CompiledValidator
	imports           Imports
	tokenizer         Tokenizer
	paramResolver     parameters.Resolver
	argResolver       ArgResolver
}

func NewCompilerOld(
	inputValidator InputValidator,
	compiledValidator CompiledValidator,
	imports Imports,
	tokenizer Tokenizer,
	paramResolver parameters.Resolver,
	argResolver ArgResolver,
) *CompilerOld {
	return &CompilerOld{
		inputValidator:    inputValidator,
		compiledValidator: compiledValidator,
		imports:           imports,
		tokenizer:         tokenizer,
		paramResolver:     paramResolver,
		argResolver:       argResolver,
	}
}

type compilerError struct {
	error
}

func throwCompilerError(err error, msg ...string) {
	if err != nil {
		if len(msg) > 0 {
			err = fmt.Errorf("%s: %s", msg[0], err.Error())
		}
		panic(compilerError{err})
	}
}

func (c CompilerOld) Compile(i input.DTO) (result compiled.DTO, err error) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}

		if cErr, ok := recovered.(compilerError); ok {
			result = compiled.DTO{}
			err = cErr
			return
		}

		panic(recovered)
	}()

	c.validateInput(i)
	c.handleMeta(i, &result)
	c.handleParams(i, &result)
	c.handleServices(i, &result)
	c.validateCompiled(result)

	return
}

func (c CompilerOld) validateInput(i input.DTO) {
	throwCompilerError(c.inputValidator.Validate(i))
}

func (c CompilerOld) validateCompiled(o compiled.DTO) {
	throwCompilerError(c.compiledValidator.Validate(o))
}

func (c CompilerOld) handleMeta(i input.DTO, result *compiled.DTO) {
	result.Meta.Pkg = i.Meta.Pkg
	result.Meta.ContainerType = i.Meta.ContainerType
	c.handleMetaImport(i.Meta.Imports)
	c.handleMetaFuncs(i.Meta.Functions)
}

func (c CompilerOld) handleMetaImport(imports map[string]string) {
	for a, p := range imports {
		throwCompilerError(
			c.imports.RegisterPrefix(a, sanitizeImport(p)),
			"cannot register alias",
		)
	}
}

var (
	regexMetaGoFn = regexp.MustCompile(`\A` + regex.MetaGoFn + `\z`)
)

func (c CompilerOld) handleMetaFuncs(funcs map[string]string) {
	for fn, goFn := range funcs {
		_, m := regex.Match(regexMetaGoFn, goFn)
		c.tokenizer.RegisterFunction(sanitizeImport(m["import"]), m["fn"], fn)
	}
}

func (c CompilerOld) handleParams(i input.DTO, result *compiled.DTO) {
	var names []string
	for n, _ := range i.Params { //nolint:gosimple
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		v := i.Params[n]
		param, err := c.paramResolver.Resolve(v)
		if err != nil {
			throwCompilerError(
				err,
				fmt.Sprintf("cannot resolve param `%s`", n),
			)
		}
		result.Params = append(
			result.Params,
			compiled.Param{
				Name:      n,
				Code:      param.Code,
				Raw:       param.Raw,
				DependsOn: param.DependsOn,
			},
		)
	}

	sort.SliceStable(result.Params, func(i, j int) bool {
		return result.Params[i].Name < result.Params[j].Name
	})
}

func (c CompilerOld) handleServices(i input.DTO, result *compiled.DTO) {
	for n, s := range i.Services {
		result.Services = append(result.Services, c.handleService(n, s))
	}
	sort.SliceStable(
		result.Services,
		func(i, j int) bool {
			return result.Services[i].Name < result.Services[j].Name
		},
	)
}

func (c CompilerOld) handleService(name string, s input.Service) compiled.Service {
	if s.Todo {
		return compiled.Service{
			Name: name,
			Todo: true,
		}
	}

	return compiled.Service{
		Name:        name,
		Getter:      s.Getter,
		Type:        c.handleServiceType(s.Type),
		Value:       c.handleServiceValue(s.Value),
		Constructor: c.handleServiceConstructor(s.Constructor),
		Args:        c.handleServiceArgs(fmt.Sprintf("service `%s`", name), s.Args),
		Calls:       c.handleServiceCalls(name, s.Calls),
		Fields:      c.handleServiceFields(name, s.Fields),
		Tags:        c.handleServiceTags(s.Tags),
		Disposable:  s.Disposable,
		Todo:        false,
	}
}

func (c CompilerOld) handleServiceType(serviceType string) string {
	_, m := regex.Match(regexServiceType, serviceType)
	t := m["type"]
	if m["import"] != "" {
		t = c.imports.GetAlias(sanitizeImport(m["import"])) + "." + t
	}
	return m["ptr"] + t
}

func (c CompilerOld) handleServiceValue(serviceValue string) string {
	if serviceValue == "" {
		return ""
	}

	_, m := regex.Match(regexServiceValue, serviceValue)

	if m["v1"] != "" {
		parts := make([]string, 0)
		if m["import"] != "" {
			parts = append(parts, c.imports.GetAlias(sanitizeImport(m["import"])))
		}
		if m["struct"] != "" {
			parts = append(parts, m["struct"]+"{}")
		}
		return strings.Join(append(parts, m["value"]), ".")
	}

	parts := make([]string, 0)
	if m["import2"] != "" {
		parts = append(parts, c.imports.GetAlias(sanitizeImport(m["import2"])))
	}
	return m["ptr2"] + strings.Join(append(parts, m["struct2"]), ".") + "{}"
}

func (c CompilerOld) handleServiceConstructor(serviceConstructor string) string {
	_, m := regex.Match(regexServiceConstructor, serviceConstructor)
	r := ""
	if m["import"] != "" {
		r = c.imports.GetAlias(sanitizeImport(m["import"])) + "."
	}
	return r + m["fn"]
}

func (c CompilerOld) handleServiceArgs(errorPrefix string, args []interface{}) (res []compiled.Arg) {
	for i, a := range args {
		arg, err := c.argResolver.Resolve(a)
		throwCompilerError(
			err,
			fmt.Sprintf("%s: cannot resolve arg%d", errorPrefix, i),
		)
		res = append(res, arg)
	}
	return
}

func (c CompilerOld) handleServiceCalls(serviceName string, calls []input.Call) (res []compiled.Call) {
	for _, raw := range calls {
		call := compiled.Call{
			Method: raw.Method,
			Args: c.handleServiceArgs(
				fmt.Sprintf("service: `%s`: call `%s`", serviceName, raw.Method),
				raw.Args,
			),
			Immutable: raw.Immutable,
		}
		res = append(res, call)
	}
	return
}

func (c CompilerOld) handleServiceFields(serviceName string, fields map[string]interface{}) (res []compiled.Field) {
	for n, f := range fields {
		arg, err := c.argResolver.Resolve(f)
		throwCompilerError(
			err,
			fmt.Sprintf("service `%s`: field `%s`", serviceName, n),
		)
		field := compiled.Field{
			Name:  n,
			Value: arg,
		}
		res = append(res, field)
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	return
}

func (c CompilerOld) handleServiceTags(tags []input.Tag) (r []compiled.Tag) {
	for _, t := range tags {
		r = append(r, compiled.Tag{
			Name:     t.Name,
			Priority: t.Priority,
		})
	}
	return
}
