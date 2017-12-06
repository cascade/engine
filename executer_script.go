package engine

import (
	"fmt"
	"os/exec"

	"github.com/cascade/protocol"
)

// ScriptExecuter ...
type ScriptExecuter struct {
	ScriptPath string
	WorkDir    string
	Env        []string
}

// Prepare ...
func (ex *ScriptExecuter) Prepare(inputs *protocol.Inputs, params *protocol.Parameters) error {
	ex.Env = []string{
		"EXAMPLE_INPUT=foobarbaz.txt",
	}
	return nil
}

// Execute ...
func (ex *ScriptExecuter) Execute() error {
	cmd := exec.Command(ex.ScriptPath)
	cmd.Env = ex.Env
	cmd.Dir = ex.WorkDir
	out, err := cmd.Output()
	if err == nil {
		fmt.Println(string(out))
	}
	return err
}

// Finalize ...
func (ex *ScriptExecuter) Finalize(outputs *protocol.Outputs) error {
	return nil
}
