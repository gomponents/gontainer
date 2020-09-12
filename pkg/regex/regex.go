package regex

import (
	"regexp"
)

// MustCompileWrapped wraps input by `\A` and `\z` and compiles into regexp.Regexp struct.
func MustCompileWrapped(r string) *regexp.Regexp {
	return regexp.MustCompile(`\A` + r + `\z`)
}

func Match(r *regexp.Regexp, s string) (bool, map[string]string) {
	if !r.MatchString(s) {
		return false, nil
	}

	match := r.FindStringSubmatch(s)
	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return true, result
}
