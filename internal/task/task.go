package task

import "time"

// Status is where a task sits in its lifecycle.
type Status string

const (
	Pending Status = "pending"
	Running Status = "running"
	Done    Status = "done"
	Failed  Status = "failed"
)

// Task is one unit of work: a job type plus the payload to run it with.
type Task struct {
	ID         string
	Type       string
	Payload    []byte
	Status     Status
	Result     string
	Error      string
	CreatedAt  time.Time
	StartedAt  time.Time
	FinishedAt time.Time
}
