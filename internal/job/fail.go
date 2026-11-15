package job

import (
	"context"
	"errors"
)

type Fail struct{}

func init() { Register("fail", Fail{}) }

// Run always fails, so retry logic in later phases has something to chew on.
func (Fail) Run(ctx context.Context, payload []byte) error {
	return errors.New("fail job always fails")
}
