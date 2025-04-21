package memory

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/vadimbarashkov/workmate-test-task/internal/entity"
	"github.com/vadimbarashkov/workmate-test-task/internal/executor"
	"github.com/vadimbarashkov/workmate-test-task/internal/handler"
	"golang.org/x/sync/semaphore"
)

type taskWrapper struct {
	ctx  context.Context
	task *entity.Task
}

type Executor struct {
	queue    chan taskWrapper
	sema     *semaphore.Weighted
	shutdown atomic.Bool
	wg       sync.WaitGroup
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

		e.wg.Add(1)
		go func(tw taskWrapper) {
			defer e.sema.Release(1)
			defer e.wg.Done()
			handler.HandleTask(tw.ctx, tw.task)
		}(tw)
	}
}

func (e *Executor) Execute(ctx context.Context, task *entity.Task) error {
	if e.shutdown.Load() {
		return executor.ErrShutdown
	}

	select {
	case e.queue <- taskWrapper{ctx: ctx, task: task}:
		return nil
	default:
		return executor.ErrQueueFull
	}
}

func (e *Executor) Shutdown(ctx context.Context) error {
	if !e.shutdown.CompareAndSwap(false, true) {
		return nil
	}

	close(e.queue)

	done := make(chan struct{})
	go func() {
		e.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil
	case <-done:
		return nil
	}
}
