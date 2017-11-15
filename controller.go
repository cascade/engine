package engine

import (
	"encoding/json"
	"io"

	"github.com/ppmp/protocol"
)

// Controller ...
type Controller struct {
	Pipeline *protocol.Pipeline `json:"pipeline" yaml:"pipeline"`
	Dryrun   bool               `json:"dryrun"   yaml:"dryrun"`
}

// New ...
func New(pipeline *protocol.Pipeline) *Controller {
	return &Controller{
		Pipeline: pipeline,
	}
}

// Handle ...
func (c *Controller) Handle(dry bool) error {
	c.Dryrun = dry

	if c.Pipeline.Steps == nil && c.Pipeline.Job != nil {
		c.Pipeline.Steps = append(c.Pipeline.Steps, c.Pipeline.Job)
		c.Pipeline.Job = nil
	}

	return nil
}

// Inspect ...
func (c *Controller) Inspect(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	return encoder.Encode(c)
}
