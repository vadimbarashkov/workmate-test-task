package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vadimbarashkov/workmate-test-task/internal/entity"
)

var (
	ErrUnknownTaskType = errors.New("unknown task type")
)

func HandleTask(ctx context.Context, task *entity.Task) {
	switch task.TaskType() {
	case entity.TypeTest:
		task.SetStatus(entity.StatusRunning)
		handleTestTask(ctx, task)
	default:
		task.SetStatus(entity.StatusFailed)
		task.SetError(ErrUnknownTaskType)
	}
}

func handleTestTask(ctx context.Context, task *entity.Task) {
	select {
	case <-ctx.Done():
		task.SetStatus(entity.StatusCanceled)
	case <-time.After(3 * time.Second): // simulate work (seconds -> minutes)
		task.SetStatus(entity.StatusCompleted)
		task.SetResult([]byte(fmt.Sprintf("Payload: %s", task.Payload())))
	}
}
