package syntax

import (
	"regexp"
	"strings"

	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexServiceValue = regexp.MustCompile(`\A` + regex.ServiceValue + `\z`)
)

type ImportAliases interface {
	GetAlias(string) string
}

func SanitizeImport(i string) string {
	return strings.Trim(i, `"`)
}

func CompileServiceValue(aliases ImportAliases, expr string) string {
	_, m := regex.Match(regexServiceValue, expr)

	if m["v1"] != "" {
		parts := make([]string, 0)
		if m["import"] != "" {
			parts = append(parts, aliases.GetAlias(SanitizeImport(m["import"])))
		}
		return m["ptr"] + strings.Join(append(parts, m["value"]), ".")
	}

	parts := make([]string, 0)
	if m["import2"] != "" {
		parts = append(parts, aliases.GetAlias(SanitizeImport(m["import2"])))
	}
	return m["ptr2"] + strings.Join(append(parts, m["struct2"]), ".") + "{}"
}
