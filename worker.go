package main

import (
	"context"
	"sync"
)

// Worker is a worker for asynchronous jobs
type Worker struct {
	ctx    context.Context
	wg     sync.WaitGroup
	cancel context.CancelFunc
	err    error
}

// NewWorker creates a new Worker with the ctx
func NewWorker(ctx context.Context) *Worker {
	return &Worker{
		ctx: ctx,
	}
}

// Start starts a background job presented by f
func (w *Worker) Start(f func(ctx context.Context) error) {
	if w.cancel != nil {
		panic("worker is already started")
	}

	ctx, cancel := context.WithCancel(w.ctx)
	w.cancel = cancel

	w.wg.Add(1)
	go func() {
		w.err = f(ctx)
		w.wg.Done()
	}()
}

// Stop stops current asynchronous jobs and returns last error occurs in the job
func (w *Worker) Stop() error {
	if w.cancel == nil {
		return nil
	}
	w.cancel()
	w.cancel = nil

	w.wg.Wait()
	return w.err
}
