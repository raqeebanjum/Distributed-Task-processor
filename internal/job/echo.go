package job

import (
	"context"
	"log/slog"
)

type Echo struct{}

func init() { Register("echo", Echo{}) }

// Run logs the payload back.
func (Echo) Run(ctx context.Context, payload []byte) error {
	slog.Info("echo", "payload", string(payload))
	return nil
}
