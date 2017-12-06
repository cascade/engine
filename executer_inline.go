package engine

import "github.com/cascade/protocol"

// InlineShellExecuter ...
type InlineShellExecuter struct {
	Inline string
}

// Prepare ...
func (ex *InlineShellExecuter) Prepare(inputs protocol.Inputs, params *protocol.Parameters) error {
	return nil
}

// Execute ...
func (ex *InlineShellExecuter) Execute() error {
	return nil
}

// Finalize ...
func (ex *InlineShellExecuter) Finalize(outputs protocol.Outputs) error {
	return nil
}
