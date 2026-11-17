package task

import (
	"errors"
	"testing"
)

func TestMemoryLifecycle(t *testing.T) {
	m := NewMemory()
	created := m.Create(Task{Type: "count", Payload: []byte(`{"max": 1}`)})

	if created.ID == "" {
		t.Fatal("Create did not assign an ID")
	}
	if created.Status != Pending {
		t.Fatalf("status = %q, want pending", created.Status)
	}
	if created.CreatedAt.IsZero() {
		t.Fatal("CreatedAt is zero")
	}

	m.Start(created.ID)
	got, ok := m.Get(created.ID)
	if !ok {
		t.Fatal("Get did not find the task")
	}
	if got.Status != Running {
		t.Fatalf("status = %q, want running", got.Status)
	}
	if got.StartedAt.IsZero() {
		t.Fatal("StartedAt is zero")
	}

	m.Finish(created.ID, nil)
	got, _ = m.Get(created.ID)
	if got.Status != Done {
		t.Fatalf("status = %q, want done", got.Status)
	}
	if got.FinishedAt.IsZero() {
		t.Fatal("FinishedAt is zero")
	}
	if got.Error != "" {
		t.Fatalf("error = %q, want empty", got.Error)
	}
}

func TestMemoryFailedTask(t *testing.T) {
	m := NewMemory()
	created := m.Create(Task{Type: "fail"})

	m.Start(created.ID)
	m.Finish(created.ID, errors.New("boom"))

	got, _ := m.Get(created.ID)
	if got.Status != Failed {
		t.Fatalf("status = %q, want failed", got.Status)
	}
	if got.Error != "boom" {
		t.Fatalf("error = %q, want boom", got.Error)
	}
	if got.FinishedAt.IsZero() {
		t.Fatal("FinishedAt is zero")
	}
}

func TestMemoryListKeepsCreationOrder(t *testing.T) {
	m := NewMemory()
	first := m.Create(Task{Type: "echo"})
	second := m.Create(Task{Type: "echo"})

	list := m.List()
	if len(list) != 2 {
		t.Fatalf("got %d tasks, want 2", len(list))
	}
	if list[0].ID != first.ID || list[1].ID != second.ID {
		t.Fatalf("list = [%s %s], want [%s %s]", list[0].ID, list[1].ID, first.ID, second.ID)
	}
}
