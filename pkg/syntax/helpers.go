package syntax

import (
	"strings"

	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexServiceValue = regex.MustCompileAz(regex.ServiceValue)
)

func SanitizeImport(i string) string {
	return strings.Trim(i, `"`)
}

// CompileServiceValue expects correct expr, validation must be done earlier
func CompileServiceValue(a imports.Aliases, expr string) string {
	_, m := regex.Match(regexServiceValue, expr)

	if m["v1"] != "" {
		parts := make([]string, 0)
		if m["import"] != "" && m["import"] != `"."` {
			parts = append(parts, a.GetAlias(SanitizeImport(m["import"])))
		}
		return m["ptr"] + strings.Join(append(parts, m["value"]), ".")
	}

	parts := make([]string, 0)
	if m["import2"] != "" && m["import2"] != `"."` {
		parts = append(parts, a.GetAlias(SanitizeImport(m["import2"])))
	}
	return m["ptr2"] + strings.Join(append(parts, m["struct2"]), ".") + "{}"
}
