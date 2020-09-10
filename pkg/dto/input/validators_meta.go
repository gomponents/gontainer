package input

import (
	"fmt"
	"regexp"

	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexpMetaPkg           = regexp.MustCompile(`\A` + regex.MetaPkg + `\z`)
	regexpMetaContainerType = regexp.MustCompile(`\A` + regex.MetaContainerType + `\z`)
	regexMetaImport         = regexp.MustCompile(`\A` + regex.MetaImport + `\z`)
	regexMetaImportAlias    = regexp.MustCompile(`\A` + regex.MetaImportAlias + `\z`)
	regexMetaFn             = regexp.MustCompile(`\A` + regex.MetaFn + `\z`)
	regexMetaGoFn           = regexp.MustCompile(`\A` + regex.MetaGoFn + `\z`)
)

// DefaultMetaValidators returns validators for DTO.Meta struct.
func DefaultMetaValidators() []func(DTO) error {
	return []func(DTO) error{
		ValidateMetaPkg,
		ValidateMetaImports,
		ValidateMetaContainerType,
		ValidateMetaFunctions,
	}
}

func ValidateMetaPkg(d DTO) error {
	if d.Meta.Pkg == "" {
		return fmt.Errorf("meta.pkg cannot be empty")
	}

	if !regexpMetaPkg.MatchString(d.Meta.Pkg) {
		return fmt.Errorf(
			"invalid meta.pkg, `%s` given",
			d.Meta.Pkg,
		)
	}

	return nil
}

func ValidateMetaContainerType(d DTO) error {
	if !regexpMetaContainerType.MatchString(d.Meta.ContainerType) {
		return fmt.Errorf(
			"invalid meta.container_type, `%s` given",
			d.Meta.ContainerType,
		)
	}
	return nil
}

func ValidateMetaImports(d DTO) error {
	for a, i := range d.Meta.Imports {
		if !regexMetaImport.MatchString(i) {
			return fmt.Errorf("invalid import `%s`", i)
		}
		if !regexMetaImportAlias.MatchString(a) {
			return fmt.Errorf("invalid alias `%s`", a)
		}
	}
	return nil
}

func ValidateMetaFunctions(d DTO) error {
	for fn, goFn := range d.Meta.Functions {
		if !regexMetaFn.MatchString(fn) {
			return fmt.Errorf("invalid function `%s`", fn)
		}

		if !regexMetaGoFn.MatchString(goFn) {
			return fmt.Errorf("invalid go function `%s`", goFn)
		}
	}
	return nil
}
