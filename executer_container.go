package engine

import "github.com/cascade/protocol"

// ContainerExecuter ...
type ContainerExecuter struct {
	ContainerName string
}

// Prepare ...
func (ex *ContainerExecuter) Prepare(inputs *protocol.Inputs, params *protocol.Parameters) error {
	return nil
}

// Execute ...
func (ex *ContainerExecuter) Execute() error {
	return nil
}

// Finalize ...
func (ex *ContainerExecuter) Finalize(outputs *protocol.Outputs) error {
	return nil
}
