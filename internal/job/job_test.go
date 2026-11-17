package job

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCountStopsWhenCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	err := Count{}.Run(ctx, []byte(`{"max": 100}`))
	elapsed := time.Since(start)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v, want context.Canceled", err)
	}
	if elapsed >= 500*time.Millisecond {
		t.Fatalf("count kept going for %s after cancel", elapsed)
	}
}

func TestGetUnknownJob(t *testing.T) {
	if _, err := Get("nope"); err == nil {
		t.Fatal("expected an error for an unknown job")
	}
}
