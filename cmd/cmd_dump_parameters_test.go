package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewDumpParamsCmd(t *testing.T) {
	newCmd := func() *cobra.Command {
		cmd := NewDumpParamsCmd()
		cmd.SilenceUsage = true
		return cmd
	}

	scenarios := []cmdScenario{
		{
			cmd:  newCmd(),
			args: "-i testdata/params.yml -l 10",
			out: `Reading files...
    testdata/params.yml

  NAME     │ PARAM                                                        
 ────────── ───────────────────────────────────────────────────────────── 
  host     │ "localhost"                                                  
  hostport │ "(...)rters".MustToString(container.MustGetParam("host")) +  
           │ "(...)rters".MustToString(":") +                             
           │ "(...)rters".MustToString(container.MustGetParam("port"))    
  port     │ int(80)                                                      
`,
		},
	}

	runCmdScenarios(t, scenarios...)
}
