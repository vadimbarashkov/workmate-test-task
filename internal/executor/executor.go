package executor

import (
	"context"
	"errors"

	"github.com/vadimbarashkov/workmate-test-task/internal/entity"
)

var (
	ErrUnknownTaskType = errors.New("unknown task type")
	ErrQueueFull       = errors.New("queue full")
)

type Executor interface {
	Execute(ctx context.Context, task *entity.Task) error
}
