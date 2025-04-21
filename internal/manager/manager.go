package manager

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vadimbarashkov/workmate-test-task/internal/entity"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type TaskManager interface {
	Create(ctx context.Context, taskType entity.TaskType, payload []byte) (uuid.UUID, error)
	Get(ctx context.Context, taskID uuid.UUID) (*entity.Task, error)
}
