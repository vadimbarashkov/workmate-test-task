package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vadimbarashkov/workmate-test-task/internal/entity"
	"github.com/vadimbarashkov/workmate-test-task/internal/executor"
	"github.com/vadimbarashkov/workmate-test-task/internal/manager"
)

type TaskManager struct {
	tasks        map[string]*entity.Task
	executor     executor.Executor
	cleanupAfter time.Duration
	mu           sync.RWMutex
}

func New(ctx context.Context, executor executor.Executor, cleanupAfter time.Duration) *TaskManager {
	tm := &TaskManager{
		tasks:        make(map[string]*entity.Task),
		executor:     executor,
		cleanupAfter: cleanupAfter,
	}
	go tm.cleanup(ctx)
	return tm
}

func (m *TaskManager) cleanup(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(m.cleanupAfter):
			var toDelete []string

			m.mu.Lock()
			for id, task := range m.tasks {
				status := task.Status()
				if status == entity.StatusCompleted || status == entity.StatusFailed || status == entity.StatusCanceled {
					toDelete = append(toDelete, id)
				}
			}

			for _, k := range toDelete {
				delete(m.tasks, k)
			}
			m.mu.Unlock()
		}
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
