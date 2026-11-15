package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"taskscheduler/internal/job"
)

const timeout = 30 * time.Second

func main() {
	name := flag.String("job", "count", "job to run")
	payload := flag.String("payload", "{}", "job payload as JSON")
	flag.Parse()

	if err := run(*name, []byte(*payload)); err != nil {
		slog.Error("job failed", "job", *name, "err", err)
		os.Exit(1)
	}
}

func run(name string, payload []byte) error {
	j, err := job.Get(name)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return j.Run(ctx, payload)
}
