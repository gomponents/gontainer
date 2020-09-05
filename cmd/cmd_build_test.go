package cmd

import (
	"strings"
	"testing"
)

func TestNewBuildCmd(t *testing.T) {
	cmd := NewBuildCmd()
	cmd.SilenceUsage = true
	assertCmd(
		t,
		cmd,
		strings.Split("-i fixtures/circular-dep.yml -o /dev/null", " "),
		`Reading files...
    fixtures/circular-dep.yml
Error: circular dependency in params: firstname -> name -> firstname
`,
		"circular dependency in params: firstname -> name -> firstname",
	)
}
