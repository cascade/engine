package engine

import (
	"fmt"
	"os"
	"sync"

	"github.com/cascade/protocol"
	"github.com/otiai10/jsonindent"
)

// Driver represents an interpriter / a driver / a handler of Cascade Component.
type Driver struct {
	// Current directory.
	Current string `json:"current"`
	// Working directory.
	WorkDir string `json:"workdir"`
	// Intermediate directory.
	IntDir string `json:"intdir"`
	// Verbose log
	Verbose bool `json:"verbose"`
}

// Run ...
func (d *Driver) Run(root *protocol.Component, param *protocol.Parameters) error {

	if d.Verbose {
		jsonindent.NewEncoder(os.Stdout).Encode(root)
		jsonindent.NewEncoder(os.Stdout).Encode(param)
		jsonindent.NewEncoder(os.Stdout).Encode(d)
	}

	defaultMachine := &protocol.Machine{Provider: "local"}

	return d.run(root, defaultMachine, param)
}

// run should be called recursively.
func (d *Driver) run(component *protocol.Component, machine *protocol.Machine, param *protocol.Parameters) error {
	if component.Machine == nil {
		component.Machine = machine
	}
	switch {
	case len(component.Steps) != 0:
		return d.steps(component, param)
	case len(component.Parallel) != 0:
		return d.parallel(component, param)
	default:
		// If this component doesn't have any nested component, handle it now.
		return d.execute(component, param)
	}
}

// steps handles child-components in series.
func (d *Driver) steps(component *protocol.Component, param *protocol.Parameters) error {
	for _, child := range component.Steps {
		if err := d.run(child, component.Machine, param); err != nil {
			return err
		}
	}
	return nil
}

// parallel handles child-components in parallel.
func (d *Driver) parallel(component *protocol.Component, param *protocol.Parameters) (err error) {
	wg := new(sync.WaitGroup)
	wg.Add(len(component.Parallel))
	for _, child := range component.Parallel {
		go func(c *protocol.Component) {
			defer wg.Done()
			if e := d.run(c, component.Machine, param); e != nil {
				err = e
			}
		}(child)
	}
	wg.Wait()
	return nil
}

// execute handles Execute
func (d *Driver) execute(component *protocol.Component, param *protocol.Parameters) error {
	if component.Execute == nil {
		return fmt.Errorf("`execute` fields cannot be empty")
	}
	executer := d.getExecuter(component)
	if executer == nil {
		return fmt.Errorf("either of `inline`, `script` or `container` must be specified in `execute` clause")
	}
	if err := executer.Prepare(component.Inputs, param); err != nil {
		return err
	}
	if err := executer.Execute(); err != nil {
		return err
	}
	if err := executer.Finalize(component.Outputs); err != nil {
		return err
	}
	return nil
}

// getExecuter ...
func (d *Driver) getExecuter(component *protocol.Component) Executer {
	execute := component.Execute
	switch {
	case execute.Inline != nil:
		return &InlineShellExecuter{}
	case execute.Script != nil:
		return &ScriptExecuter{}
	case execute.Container != nil:
		return &ContainerExecuter{}
	}
	return nil
}
