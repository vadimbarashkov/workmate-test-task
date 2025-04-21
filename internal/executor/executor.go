package executor

import (
	"context"
	"errors"

	"github.com/vadimbarashkov/workmate-test-task/internal/entity"
)

var (
	ErrQueueFull = errors.New("queue full")
)

type Executor interface {
	Execute(ctx context.Context, task *entity.Task) error
}
