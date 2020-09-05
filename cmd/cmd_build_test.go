package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewBuildCmd(t *testing.T) {
	newCmd := func() *cobra.Command {
		cmd := NewBuildCmd()
		cmd.SilenceUsage = true
		return cmd
	}

	assertCmd(
		t,
		newCmd(),
		strings.Split("-i testdata/circular-dep-params.yml -o /dev/null", " "),
		`Reading files...
    testdata/circular-dep-params.yml
Error: circular dependency in params: firstname -> name -> firstname
`,
		"circular dependency in params: firstname -> name -> firstname",
	)

	assertCmd(
		t,
		newCmd(),
		strings.Split("-i testdata/circular-dep-services.yml -o /dev/null", " "),
		`Reading files...
    testdata/circular-dep-services.yml
Error: circular dependency in services: db -> storage -> db
`,
		"circular dependency in services: db -> storage -> db",
	)
}
