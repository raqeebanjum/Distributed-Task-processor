package task

import (
	"fmt"
	"sync"
	"time"
)

// Memory keeps tasks in a map guarded by a mutex. Phase 4 turns this into a Store.
type Memory struct {
	mu    sync.Mutex
	tasks map[string]*Task
	order []string
	next  int
}

func NewMemory() *Memory {
	return &Memory{tasks: map[string]*Task{}}
}

// Create stores a new pending task, assigning its ID and creation time.
func (m *Memory) Create(t Task) Task {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.next++
	t.ID = fmt.Sprintf("task-%d", m.next)
	t.Status = Pending
	t.CreatedAt = time.Now()

	m.tasks[t.ID] = &t
	m.order = append(m.order, t.ID)
	return t
}

// Get returns a copy of the task, so callers never share the stored struct.
func (m *Memory) Get(id string) (Task, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.tasks[id]
	if !ok {
		return Task{}, false
	}
	return *t, true
}

// List returns copies of every task, oldest first.
func (m *Memory) List() []Task {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := make([]Task, 0, len(m.order))
	for _, id := range m.order {
		out = append(out, *m.tasks[id])
	}
	return out
}

// Start marks a task as running and records when it began.
func (m *Memory) Start(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if t, ok := m.tasks[id]; ok {
		t.Status = Running
		t.StartedAt = time.Now()
	}
}

// Finish marks a task done, or failed with the error message.
func (m *Memory) Finish(id string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.tasks[id]
	if !ok {
		return
	}
	t.FinishedAt = time.Now()
	if err != nil {
		t.Status = Failed
		t.Error = err.Error()
		return
	}
	t.Status = Done
}
