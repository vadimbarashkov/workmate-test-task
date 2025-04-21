package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/vadimbarashkov/workmate-test-task/internal/entity"
)

func HandleTask(ctx context.Context, task *entity.Task) {
	switch task.TaskType() {
	case entity.TypeTest:
		task.SetStatus(entity.StatusRunning)
		handleTestTask(ctx, task)
	default:
		task.SetStatus(entity.StatusFailed)
		task.SetError(entity.ErrInvalidTaskType)
	}
}

func handleTestTask(ctx context.Context, task *entity.Task) {
	if len(task.Payload()) == 0 {
		task.SetStatus(entity.StatusFailed)
		task.SetError(entity.ErrMissingPayload)
		return
	}

	select {
	case <-ctx.Done():
		task.SetStatus(entity.StatusCanceled)
	case <-time.After(15 * time.Second): // simulate work (seconds -> minutes)
		task.SetStatus(entity.StatusCompleted)
		task.SetResult([]byte(fmt.Sprintf("Payload: %s", task.Payload())))
	}
}
