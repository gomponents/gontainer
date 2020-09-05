package cmd

import (
	"bytes"
	"io/ioutil"
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
