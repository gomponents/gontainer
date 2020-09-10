package input

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMetaPkg(t *testing.T) {
	scenarios := []struct {
		pkg   string
		error string
	}{
		{
			pkg:   "",
			error: "meta.pkg cannot be empty",
		},
		{
			pkg: "main",
		},
		{
			pkg:   "123",
			error: "invalid meta.pkg, `123` given",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			dto := DTO{}
			dto.Meta.Pkg = s.pkg
			err := ValidateMetaPkg(dto)
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, s.error)
		})
	}
}

func TestValidateMetaContainerType(t *testing.T) {
	scenarios := []struct {
		containerType string
		error         string
	}{
		{
			containerType: "",
			error:         "invalid meta.container_type, `` given",
		},
		{
			containerType: "myContainer123",
		},
		{
			containerType: "0MyContainer",
			error:         "invalid meta.container_type, `0MyContainer` given",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			dto := DTO{}
			dto.Meta.ContainerType = s.containerType
			err := ValidateMetaContainerType(dto)
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, s.error)
		})
	}
}

func TestValidateMetaImports(t *testing.T) {
	scenarios := []struct {
		import_ string
		alias   string
		error   string
	}{
		{
			import_: "github.com/stretchr/testify/assert",
			alias:   "assert",
		},
		{
			import_: "github.com/stretchr/testify/assert/",
			alias:   "assert",
			error:   "invalid import `github.com/stretchr/testify/assert/`",
		},
		{
			import_: "oneTwoThree",
			alias:   "$123",
			error:   "invalid alias `$123`",
		},
		{
			import_: "!!!",
			alias:   "alias",
			error:   "invalid import `!!!`",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			d := DTO{}
			d.Meta.Imports = map[string]string{
				s.alias: s.import_,
			}

			err := ValidateMetaImports(d)

			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, s.error, s)
		})
	}
}

func TestValidateMetaFunctions(t *testing.T) {
	scenarios := []struct {
		alias string
		goFn  string
		error string
	}{
		{
			alias: "env",
			goFn:  "os.Getenv",
		},
		{
			alias: "env",
			goFn:  "env",
		},
		{
			alias: "$fn",
			goFn:  "os.Getenv",
			error: "invalid function `$fn`",
		},
		{
			alias: "env",
			goFn:  "os.1Getenv",
			error: "invalid go function `os.1Getenv`",
		},
	}
	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			d := DTO{}
			d.Meta.Functions = map[string]string{
				s.alias: s.goFn,
			}
			err := ValidateMetaFunctions(d)

			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, s.error)
		})
	}
}
