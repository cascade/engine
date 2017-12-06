package engine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	execute := component.Execute
	if execute == nil {
		return fmt.Errorf("`execute` fields cannot be empty")
	}
	switch {
	case execute.Inline != nil:
		return d.executeInlineShell(component, param)
	case execute.Script != nil:
		return d.executeScript(component, param)
	case execute.Container != nil:
		return d.executeContainer(component, param)
	}
	return fmt.Errorf("either of `inline`, `script` or `container` must be specified in `execute` clause")
}

func (d *Driver) executeInlineShell(component *protocol.Component, param *protocol.Parameters) error {
	// FIXME: Only bash?
	cmd := exec.Command("bash", "-c", *component.Execute.Inline)
	cmd.Dir = d.WorkDir
	out, err := cmd.Output()
	// FIXME: hard coding
	fmt.Println("DEBUG:", string(out))
	return err
}

// executeCommandLine ...
func (d *Driver) executeScript(component *protocol.Component, param *protocol.Parameters) error {
	if component.Runtime == nil {
		scriptpath, err := filepath.Abs(*component.Execute.Script)
		if err != nil {
			return err
		}
		component.Execute.Script = &scriptpath
		cmd := exec.Command(*component.Execute.Script)
		cmd.Dir = d.WorkDir
		out, err := cmd.Output()
		// FIXME: hard coding
		fmt.Println("DEBUG:", string(out))
		return err
	}
	// Create a machine with specified runtime, and execute the script to that container.
	return nil
}

func (d *Driver) executeContainer(component *protocol.Component, param *protocol.Parameters) error {
	return nil
}
