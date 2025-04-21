package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/vadimbarashkov/workmate-test-task/internal/entity"
	"github.com/vadimbarashkov/workmate-test-task/internal/executor"
	"github.com/vadimbarashkov/workmate-test-task/internal/manager"
)

type TaskManager struct {
	tasks    map[string]*entity.Task
	executor executor.Executor
	mu       sync.RWMutex
}

func New(executor executor.Executor) *TaskManager {
	return &TaskManager{
		tasks:    make(map[string]*entity.Task),
		executor: executor,
	}
}

func (m *TaskManager) Create(ctx context.Context, taskType entity.TaskType, payload []byte) (uuid.UUID, error) {
	task := entity.NewTask(taskType, payload)

	if err := m.executor.Execute(ctx, task); err != nil {
		return uuid.Nil, fmt.Errorf("execute task: %w", err)
	}

	m.mu.Lock()
	m.tasks[task.ID().String()] = task
	m.mu.Unlock()

	return task.ID(), nil
}

func (m *TaskManager) Get(_ context.Context, taskID uuid.UUID) (*entity.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, exists := m.tasks[taskID.String()]
	if !exists {
		return nil, manager.ErrTaskNotFound
	}

	return task, nil
}

// TODO: periodically delete executed, failed or canceled tasks
