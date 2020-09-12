package input

import (
	"fmt"
	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexpMetaPkg                  = regex.MustCompileWrapped(regex.MetaPkg)
	regexpMetaContainerType        = regex.MustCompileWrapped(regex.MetaContainerType)
	regexpMetaContainerConstructor = regex.MustCompileWrapped(regex.MetaContainerConstructor)
	regexMetaImport                = regex.MustCompileWrapped(regex.MetaImport)
	regexMetaImportAlias           = regex.MustCompileWrapped(regex.MetaImportAlias)
	regexMetaFn                    = regex.MustCompileWrapped(regex.MetaFn)
	regexMetaGoFn                  = regex.MustCompileWrapped(regex.MetaGoFn)
)

// DefaultMetaValidators returns validators for DTO.Meta struct.
func DefaultMetaValidators() []func(DTO) error {
	return []func(DTO) error{
		ValidateMetaPkg,
		ValidateMetaImports,
		ValidateMetaContainerType,
		ValidateMetaContainerConstructor,
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

func ValidateMetaContainerConstructor(d DTO) error {
	if !regexpMetaContainerConstructor.MatchString(d.Meta.ContainerConstructor) {
		return fmt.Errorf(
			"invalid meta.container_constructor, `%s` given",
			d.Meta.ContainerConstructor,
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
