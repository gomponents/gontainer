package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
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
		{
			cmd:  newCmd(),
			args: "-i testdata/empty.yml",
			out: `Reading files...
    testdata/empty.yml
Could not find any parameters
`,
			error: "",
		},
		{
			cmd:  newCmd(),
			args: "-i testdata/params.yml -l 5", // minimum length is 10
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
		{
			cmd:  newCmd(),
			args: "-i testdata/empty.yml",
			out: `Reading files...
    testdata/empty.yml
Could not find any parameters
`,
			error: "",
		},
	}

	runCmdScenarios(t, scenarios...)
}

func Test_fakeImports_GetImports(t *testing.T) {
	assert.Nil(t, fakeImports{}.GetImports())
}

func Test_fakeImports_RegisterPrefix(t *testing.T) {
	assert.NoError(t, fakeImports{}.RegisterPrefix("pkg", "my/repo/pkg"))
}
