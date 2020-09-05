package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewDumpParamsCmd(t *testing.T) {
	newCmd := func() *cobra.Command {
		cmd := NewDumpParamsCmd()
		cmd.SilenceUsage = true
		return cmd
	}

	assertCmd(
		t,
		newCmd(),
		strings.Split("-i testdata/params.yml", " "),
		`Reading files...
    testdata/params.yml

  NAME     │ PARAM                                                           
 ────────── ──────────────────────────────────────────────────────────────── 
  host     │ "localhost"                                                     
  hostport │ "(...)/exporters".MustToString(container.MustGetParam("host"))  
           │ + "(...)/exporters".MustToString(":") +                         
           │ "(...)/exporters".MustToString(container.MustGetParam("port"))  
  port     │ int(80)                                                         
`,
		"",
	)
}
