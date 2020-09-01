package imports

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	regexNoAlphaNum = regexp.MustCompile("[^a-zA-Z0-9]")
)

type Import struct {
	Path  string
	Alias string
}

type Imports interface {
	GetAlias(string) string
	GetImports() []Import
	RegisterPrefix(shortcut string, path string) error
}

type SimpleImports struct {
	counter      int64
	imports      map[string]Import
	importsSlice []Import
	shortcuts    map[string]string
}

func (s *SimpleImports) RegisterPrefix(shortcut string, path string) error {
	if _, ok := s.shortcuts[shortcut]; ok {
		return fmt.Errorf("shortcut `%s` is already registered", shortcut)
	}

	s.shortcuts[shortcut] = path
	return nil
}

func (s *SimpleImports) GetAlias(path string) string {
	path = s.decorateImport(path)

	if i, ok := s.imports[path]; ok {
		return i.Alias
	}

	parts := strings.Split(path, "/")

	alias := parts[len(parts)-1]
	if len(parts) >= 2 {
		alias = parts[len(parts)-2] + "_" + alias
	}
	alias = regexNoAlphaNum.ReplaceAllString(alias, "_")
	alias = fmt.Sprintf("i%s_%s", strconv.FormatInt(s.counter, 16), alias)

	i := Import{
		Path:  path,
		Alias: alias,
	}
	s.imports[path] = i
	s.counter++
	s.importsSlice = append(s.importsSlice, i)

	return i.Alias
}

func (s *SimpleImports) GetImports() []Import {
	//nolint:gosimple,staticcheck
	return append(s.importsSlice)
}

func (s *SimpleImports) decorateImport(i string) string {
	for shortcut, path := range s.shortcuts {
		if strings.Index(i, shortcut) == 0 {
			return strings.Replace(i, shortcut, path, 1)
		}
	}

	return i
}

func NewSimpleImports() *SimpleImports {
	return &SimpleImports{
		imports:   make(map[string]Import),
		shortcuts: make(map[string]string),
	}
}
