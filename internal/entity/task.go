package entity

import "github.com/google/uuid"

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
	default:
		return "unknown"
	}
}

type Task struct {
	ID      uuid.UUID
	Type    TaskType
	Status  TaskStatus
	Payload []byte
	Result  []byte
	Error   error
}

func NewTask(taskType TaskType, payload []byte) *Task {
	return &Task{
		ID:     uuid.New(),
		Type:   taskType,
		Status: StatusPending,
	}
}
