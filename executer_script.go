package engine

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cascade/protocol"
)

// ScriptExecuter ...
type ScriptExecuter struct {
	ScriptPath string
	WorkDir    string
	Env        []string
}

// Prepare ...
func (ex *ScriptExecuter) Prepare(inputs protocol.Inputs, params *protocol.Parameters) error {
	for key, required := range inputs {
		provided, ok := params.Inputs[key]
		if !ok {
			return fmt.Errorf("input `%s` is required but not provided", key)
		}
		if provided.Type != required.Type {
			return fmt.Errorf("required type for `%s` is `%s`, but `%s` is provided", key, required.Type, provided.Type)
		}
		switch provided.Type {
		case "file":
			destpath, err := ex.prepareFile(required, provided)
			if err != nil {
				return err
			}
			ex.Env = append(ex.Env, fmt.Sprintf("%s=%s", key, destpath))
		}
		fmt.Printf("%+v\n", required)
		fmt.Printf("%+v\n", provided)
	}
	return nil
}

// prepareFile
func (ex *ScriptExecuter) prepareFile(required, provided *protocol.Input) (string, error) {
	destpath := filepath.Join(ex.WorkDir, filepath.Base(provided.URL))
	f, err := os.Create(destpath)
	if err != nil {
		return destpath, err
	}
	defer f.Close()
	res, err := http.Get(provided.URL)
	if err != nil {
		return destpath, err
	}
	defer res.Body.Close()
	if _, err := io.Copy(f, res.Body); err != nil {
		return destpath, err
	}
	return destpath, nil
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
func (ex *ScriptExecuter) Finalize(outputs protocol.Outputs) error {
	return nil
}
