package job

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

type Count struct{}

func init() { Register("count", Count{}) }

// Run logs the numbers 1 to max, pausing between each one, until it finishes
// or the context is cancelled.
func (Count) Run(ctx context.Context, payload []byte) error {
	opts := struct {
		Max int `json:"max"`
	}{Max: 10}
	if len(payload) > 0 {
		if err := json.Unmarshal(payload, &opts); err != nil {
			return fmt.Errorf("decode count payload: %w", err)
		}
	}

	for i := 1; i <= opts.Max; i++ {
		slog.Info("count", "n", i)
		if i == opts.Max {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
	return nil
}
