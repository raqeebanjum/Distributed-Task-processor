package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"taskscheduler/internal/task"
	"taskscheduler/internal/worker"
)

func main() {
	size := flag.Int("pool", 3, "number of workers")
	flag.Parse()

	store := task.NewMemory()
	pool := worker.New(store, *size)
	pool.Start(context.Background())

	for i := 0; i < 5; i++ {
		_, err := pool.Submit(task.Task{
			Type:    "count",
			Payload: []byte(`{"max": 3}`),
		})
		if err != nil {
			slog.Error("submit failed", "err", err)
			os.Exit(1)
		}
	}

	pool.Wait()
	pool.Stop()

	fmt.Printf("%-7s %-6s %-7s %-13s %-13s %s\n", "id", "type", "status", "started", "finished", "error")
	for _, t := range store.List() {
		fmt.Printf("%-7s %-6s %-7s %-13s %-13s %s\n", t.ID, t.Type, t.Status,
			t.StartedAt.Format("15:04:05.000"), t.FinishedAt.Format("15:04:05.000"), t.Error)
	}
}
