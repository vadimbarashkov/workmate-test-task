package executor

import (
	"context"
	"errors"

	"github.com/vadimbarashkov/workmate-test-task/internal/entity"
)

var (
	ErrQueueFull = errors.New("queue full")
	ErrShutdown  = errors.New("shutdown")
)

type Executor interface {
	Execute(ctx context.Context, task *entity.Task) error
	Shutdown(ctx context.Context) error
}
