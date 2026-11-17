package worker

import (
	"context"
	"log/slog"
	"sync"

	"taskscheduler/internal/job"
	"taskscheduler/internal/task"
)

// Pool runs tasks on a fixed number of worker goroutines.
type Pool struct {
	store *task.Memory
	queue chan task.Task
	size  int
	tasks sync.WaitGroup // submitted tasks that have not finished
	live  sync.WaitGroup // worker goroutines
}

func New(store *task.Memory, size int) *Pool {
	return &Pool{
		store: store,
		queue: make(chan task.Task),
		size:  size,
	}
}

// Start launches the workers. They run until Stop is called.
func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.size; i++ {
		p.live.Add(1)
		go p.work(ctx)
	}
}

// Submit stores a task as pending and queues it, returning the stored task.
func (p *Pool) Submit(t task.Task) (task.Task, error) {
	if _, err := job.Get(t.Type); err != nil {
		return task.Task{}, err
	}

	created := p.store.Create(t)
	p.tasks.Add(1)
	p.queue <- created
	return created, nil
}

// Wait blocks until every submitted task has finished.
func (p *Pool) Wait() {
	p.tasks.Wait()
}

// Stop closes the queue and waits for the workers to exit.
func (p *Pool) Stop() {
	close(p.queue)
	p.live.Wait()
}

func (p *Pool) work(ctx context.Context) {
	defer p.live.Done()

	for t := range p.queue {
		p.run(ctx, t)
		p.tasks.Done()
	}
}

func (p *Pool) run(ctx context.Context, t task.Task) {
	slog.Info("task started", "id", t.ID, "type", t.Type)
	p.store.Start(t.ID)

	j, err := job.Get(t.Type)
	if err == nil {
		err = j.Run(ctx, t.Payload)
	}
	p.store.Finish(t.ID, err)

	if err != nil {
		slog.Warn("task failed", "id", t.ID, "err", err)
		return
	}
	slog.Info("task done", "id", t.ID)
}
