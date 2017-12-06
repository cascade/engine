package engine

import "github.com/cascade/protocol"

// ScriptExecuter ...
type ScriptExecuter struct {
	ScriptPath string
}

// Prepare ...
func (ex *ScriptExecuter) Prepare(inputs *protocol.Inputs, params *protocol.Parameters) error {
	return nil
}

// Execute ...
func (ex *ScriptExecuter) Execute() error {
	return nil
}

// Finalize ...
func (ex *ScriptExecuter) Finalize(outputs *protocol.Outputs) error {
	return nil
}
