package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/vadimbarashkov/workmate-test-task/internal/entity"
	"github.com/vadimbarashkov/workmate-test-task/internal/manager"
)

type TaskManager struct {
	mu    sync.RWMutex
	tasks map[string]*entity.Task
}

func New() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*entity.Task),
	}
}

func (m *TaskManager) Create(_ context.Context, taskType entity.TaskType, payload []byte) (uuid.UUID, error) {
	task := entity.NewTask(taskType, payload)
	m.mu.Lock()
	m.tasks[task.ID().String()] = task
	m.mu.Unlock()
	// TODO: add to queue
	return task.ID(), nil
}

func (m *TaskManager) Get(_ context.Context, taskID uuid.UUID) (*entity.Task, error) {
	m.mu.RLock()
	task, exists := m.tasks[taskID.String()]
	if !exists {
		return nil, manager.ErrTaskNotFound
	}
	m.mu.RUnlock()
	return task, nil
}
