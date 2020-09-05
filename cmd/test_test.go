package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func assertCmd(t *testing.T, cmd *cobra.Command, args []string, out string, err string) {
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs(args)

	cmdErr := cmd.Execute()
	if err == "" {
		if !assert.NoError(t, cmdErr) {
			return
		}
	} else {
		if !assert.EqualError(t, cmdErr, err) {
			return
		}
	}

	cmdOut, ioErr := ioutil.ReadAll(b)
	if ioErr != nil {
		t.Fatal(ioErr)
		return
	}

	assert.Equal(t, out, string(cmdOut))
}

type cmdScenario struct {
	cmd   *cobra.Command
	args  string
	out   string
	error string
}

func runCmdScenarios(t *testing.T, scenarios ...cmdScenario) {
	for i, s := range scenarios {
		var args []string
		if s.args != "" {
			args = strings.Split(s.args, " ")
		}
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			assertCmd(
				t,
				s.cmd,
				args,
				s.out,
				s.error,
			)
		})
	}
}
