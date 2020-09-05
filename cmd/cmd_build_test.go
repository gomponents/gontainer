package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewBuildCmd(t *testing.T) {
	newCmd := func() *cobra.Command {
		cmd := NewBuildCmd()
		cmd.SilenceUsage = true
		return cmd
	}

	scenarios := []cmdScenario{
		{
			cmd:  newCmd(),
			args: "-i testdata/circular-dep-params.yml -o /dev/null",
			out: `Reading files...
    testdata/circular-dep-params.yml
Error: circular dependency in params: firstname -> name -> firstname
`,
			error: "circular dependency in params: firstname -> name -> firstname",
		},
		{
			cmd:  newCmd(),
			args: "-i testdata/circular-dep-services.yml -o /dev/null",
			out: `Reading files...
    testdata/circular-dep-services.yml
Error: circular dependency in services: db -> storage -> db
`,
			error: "circular dependency in services: db -> storage -> db",
		},
		{
			cmd:   newCmd(),
			args:  "-i [] -o /dev/null",
			out:   "Error: pattern: `[]`: syntax error in pattern\n",
			error: "pattern: `[]`: syntax error in pattern",
		},
		{
			cmd:   newCmd(),
			args:  "-i foo/bar/*.yml -o /dev/null",
			out:   "Error: cannot find any configuration file\n",
			error: "cannot find any configuration file",
		},
		{
			cmd:  newCmd(),
			args: "-i ../examples/library/container/gontainer.yml -o /dev/null",
			out: "Reading files...\n" +
				"    ../examples/library/container/gontainer.yml\n" +
				"Successfully generated container to file: `/dev/null`\n",
			error: "",
		},
		{
			cmd:  newCmd(),
			args: "-i ../examples/library/container/gontainer.yml -o /",
			out: `Reading files...
    ../examples/library/container/gontainer.yml
Error: open /: is a directory
`,
			error: "open /: is a directory",
		},
		{
			cmd:  newCmd(),
			args: "-i / -o /dev/null",
			out: "Reading files...\n" +
				"    /\n" +
				"Error: error has occurred during opening file `/`: read /: is a directory\n",
			error: "error has occurred during opening file `/`: read /: is a directory",
		},
		{
			cmd:  newCmd(),
			args: "-i testdata/invalid.yml -o /dev/null",
			out: "Reading files...\n" +
				"    testdata/invalid.yml\n" +
				"Error: error has occurred during parsing yaml file `testdata/invalid.yml`: yaml: line 2: did not find expected node content\n",
			error: "error has occurred during parsing yaml file `testdata/invalid.yml`: yaml: line 2: did not find expected node content",
		},
		{ // todo do not print regex in cli
			cmd:  newCmd(),
			args: "-i testdata/invalid-constructor.yml -o /dev/null",
			out: "Reading files...\n" +
				"    testdata/invalid-constructor.yml\n" +
				"Error: service `db`: constructor must match `\\A((?P<import>(([A-Za-z](\\/?[A-Z-a-z0-9._-])*)|(\"[A-Za-z](\\/?[A-Z-a-z0-9._-])*\")|\"\\.\"))\\.)?(?P<fn>[A-Za-z][A-Za-z0-9_]*)\\z`, `New DB` given\n",
			error: "service `db`: constructor must match `\\A((?P<import>(([A-Za-z](\\/?[A-Z-a-z0-9._-])*)|(\"[A-Za-z](\\/?[A-Z-a-z0-9._-])*\")|\"\\.\"))\\.)?(?P<fn>[A-Za-z][A-Za-z0-9_]*)\\z`, `New DB` given",
		},
	}

	runCmdScenarios(t, scenarios...)
}
