package engine

import "github.com/cascade/protocol"

// Executer ...
type Executer interface {
	Prepare(protocol.Inputs, *protocol.Parameters) error
	Execute() error
	Finalize(protocol.Outputs) error
}
