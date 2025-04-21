package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/vadimbarashkov/workmate-test-task/internal/entity"
	"github.com/vadimbarashkov/workmate-test-task/internal/executor"
	"golang.org/x/sync/semaphore"
)

type taskWrapper struct {
	ctx  context.Context
	task *entity.Task
}

type Executor struct {
	queue chan taskWrapper
	sema  *semaphore.Weighted
}

func New(queueSize int, maxWorkers int64) *Executor {
	if queueSize <= 0 {
		panic("queue size must be positive")
	}
	if maxWorkers <= 0 {
		panic("max workers must be positive")
	}

	e := &Executor{
		queue: make(chan taskWrapper, queueSize),
		sema:  semaphore.NewWeighted(maxWorkers),
	}

	go e.dispatch()
	return e
}

func (e *Executor) dispatch() {
	for tw := range e.queue {
		if err := e.sema.Acquire(tw.ctx, 1); err != nil {
			tw.task.SetStatus(entity.StatusCanceled)
			continue
		}

		go func(tw taskWrapper) {
			defer e.sema.Release(1)
			e.handleTask(tw.ctx, tw.task)
		}(tw)
	}
}

func (e *Executor) handleTask(ctx context.Context, task *entity.Task) {
	switch task.TaskType() {
	case entity.TypeTest:
		e.handleTestTask(ctx, task)
	default:
		task.SetStatus(entity.StatusFailed)
		task.SetError(executor.ErrUnknownTaskType)
	}
}

func (e *Executor) handleTestTask(ctx context.Context, task *entity.Task) {
	select {
	case <-ctx.Done():
		task.SetStatus(entity.StatusCanceled)
	case <-time.After(3 * time.Second): // simulate work (seconds -> minutes)
		task.SetStatus(entity.StatusCompleted)
		task.SetResult([]byte(fmt.Sprintf("Payload: %s", task.Payload())))
	}
}

func (e *Executor) Execute(ctx context.Context, task *entity.Task) error {
	select {
	case e.queue <- taskWrapper{ctx: ctx, task: task}:
		task.SetStatus(entity.StatusRunning)
		return nil
	default:
		return executor.ErrQueueFull
	}
}
