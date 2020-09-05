package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"gopkg.in/yaml.v2"
)

type readConfig struct {
	printer
	patterns []string
}

func newReadConfig(w io.Writer, patterns []string) *readConfig {
	return &readConfig{printer: printer{w: w}, patterns: patterns}
}

func (r readConfig) run(i *input.DTO, _ *compiled.DTO) error {
	files, fErr := findFiles(r.patterns)
	if fErr != nil {
		return fErr
	}

	if len(files) == 0 {
		return fmt.Errorf("cannot find any configuration file")
	}

	r.printf("Reading files...\n")
	for _, f := range files {
		r.printf("%s%s\n", indent, f)
		if file, err := ioutil.ReadFile(f); err != nil {
			return fmt.Errorf("error has occurred during opening file `%s`: %s", f, err.Error())
		} else {
			if yamlErr := yaml.Unmarshal(file, &i); yamlErr != nil {
				return fmt.Errorf("error has occurred during parsing yaml file `%s`: %s", f, yamlErr.Error())
			}
		}
	}

	return nil
}

// findFiles returns list of files found by given patterns.
// Files are sorted by index of pattern then name.
func findFiles(patterns []string) ([]string, error) {
	result := make([]string, 0)
	for _, p := range patterns {
		matches, err := filepath.Glob(p)
		if err != nil {
			return nil, fmt.Errorf("pattern: `%s`: %s", p, err.Error())
		}
		for i, m := range matches {
			matches[i] = filepath.Clean(m)
		}
		sort.Strings(matches)
		result = append(result, matches...)
	}

	return result, nil
}
