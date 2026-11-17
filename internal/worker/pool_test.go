package worker

import (
	"context"
	"sync"
	"testing"
	"time"

	"taskscheduler/internal/job"
	"taskscheduler/internal/task"
)

// barrier blocks until all tasks reach Run at the same time, so the test fails
// instead of passing if the pool runs tasks one after another.
type barrier struct{ wg *sync.WaitGroup }

func (b barrier) Run(ctx context.Context, payload []byte) error {
	b.wg.Done()
	b.wg.Wait()
	return nil
}

func TestPoolRunsTasksConcurrently(t *testing.T) {
	var running sync.WaitGroup
	running.Add(3)
	job.Register("test-barrier", barrier{wg: &running})

	store := task.NewMemory()
	pool := New(store, 3)
	pool.Start(context.Background())
	defer pool.Stop()

	for i := 0; i < 3; i++ {
		if _, err := pool.Submit(task.Task{Type: "test-barrier"}); err != nil {
			t.Fatal(err)
		}
	}

	done := make(chan struct{})
	go func() {
		pool.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("tasks never ran at the same time")
	}

	for _, got := range store.List() {
		if got.Status != task.Done {
			t.Fatalf("%s status = %q, want done", got.ID, got.Status)
		}
	}
}

func TestPoolMarksFailedTask(t *testing.T) {
	store := task.NewMemory()
	pool := New(store, 1)
	pool.Start(context.Background())
	defer pool.Stop()

	if _, err := pool.Submit(task.Task{Type: "fail"}); err != nil {
		t.Fatal(err)
	}
	pool.Wait()

	got, ok := store.Get("task-1")
	if !ok {
		t.Fatal("task was not stored")
	}
	if got.Status != task.Failed {
		t.Fatalf("status = %q, want failed", got.Status)
	}
	if got.Error == "" {
		t.Fatal("error message is empty")
	}
	if got.StartedAt.IsZero() || got.FinishedAt.IsZero() {
		t.Fatal("timestamps are not set")
	}
}

func TestSubmitUnknownJob(t *testing.T) {
	pool := New(task.NewMemory(), 1)
	pool.Start(context.Background())
	defer pool.Stop()

	if _, err := pool.Submit(task.Task{Type: "nope"}); err == nil {
		t.Fatal("expected an error for an unknown job type")
	}
}
