package entity

import (
	"sync"

	"github.com/google/uuid"
)

type TaskType int

const (
	TypeTest TaskType = iota
)

func (t TaskType) String() string {
	switch t {
	case TypeTest:
		return "test"
	default:
		return "unknown"
	}
}

type TaskStatus int

const (
	StatusPending TaskStatus = iota
	StatusRunning
	StatusCompleted
	StatusFailed
	StatusCanceled
)

func (s TaskStatus) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusRunning:
		return "running"
	case StatusCompleted:
		return "completed"
	case StatusFailed:
		return "failed"
	case StatusCanceled:
		return "canceled"
	default:
		return "unknown"
	}
}

type Task struct {
	id       uuid.UUID
	taskType TaskType
	status   TaskStatus
	payload  []byte
	result   []byte
	error    error
	mu       sync.RWMutex
}

func NewTask(taskType TaskType, payload []byte) *Task {
	return &Task{
		id:       uuid.New(),
		taskType: taskType,
		status:   StatusPending,
		payload:  payload,
	}
}

func (t *Task) ID() uuid.UUID {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.id
}

func (t *Task) TaskType() TaskType {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.taskType
}

func (t *Task) Status() TaskStatus {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.status
}

func (t *Task) SetStatus(status TaskStatus) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.status = status
}

func (t *Task) Payload() []byte {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.payload
}

func (t *Task) Result() []byte {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.result
}

func (t *Task) SetResult(result []byte) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.result = result
}

func (t *Task) Error() error {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.error
}

func (t *Task) SetError(err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.error = err
}
