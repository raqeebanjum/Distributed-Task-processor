package job

import (
	"context"
	"fmt"
)

// Job is a unit of work the runner can execute.
type Job interface {
	Run(ctx context.Context, payload []byte) error
}

var registry = map[string]Job{}

// Register makes a job available by name. Jobs call this from init.
func Register(name string, j Job) {
	registry[name] = j
}

// Get looks up a job by name.
func Get(name string) (Job, error) {
	j, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown job %q", name)
	}
	return j, nil
}
