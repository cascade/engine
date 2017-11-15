package engine

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"github.com/ppmp/protocol"
)

// Controller ...
type Controller struct {
	Pipeline *protocol.Pipeline
}

// ParseFile ...
func ParseFile(f *os.File) (*Controller, error) {
	pipeline := new(protocol.Pipeline)
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	ext := filepath.Ext(f.Name())
	if err := getUnmarshalFunc(ext)(b, pipeline); err != nil {
		return nil, err
	}
	return New(pipeline), nil
}

// getUnmarshalFunc ...
func getUnmarshalFunc(ext string) func(b []byte, v interface{}) error {
	switch ext {
	case "json":
		return json.Unmarshal
	case "yaml", "yml":
		return yaml.Unmarshal
	default:
		return yaml.Unmarshal
	}
}

// New ...
func New(pipeline *protocol.Pipeline) *Controller {
	return &Controller{
		Pipeline: pipeline,
	}
}

// Handle ...
func Handle() error {
	return nil
}
