package worker

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNoWorkers   = errors.New("no workers provided")
	ErrNoExecuters = errors.New("no executers provided")
)

type pool struct {
	workers   []Worker
	executers []Executer
}

func NewPool(options ...optFunc) (*pool, error) {
	w := pool{}

	for _, opt := range options {
		opt(&w)
	}

	if w.workers == nil {
		return nil, ErrNoWorkers
	}

	if w.executers == nil {
		return nil, ErrNoExecuters
	}

	return &w, nil
}

func (w pool) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	workerChannel := make(chan Worker, len(w.workers))

	for _, e := range w.executers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for w := range workerChannel {
				select {
				case <-ctx.Done():
					return
				default:
					_ = e.Execute(ctx, w.Do)
				}
			}
		}()

	}

	for _, w := range w.workers {
		workerChannel <- w
	}
	close(workerChannel)

	wg.Wait()

	if err := ctx.Err(); err != nil {
		return err
	}

	return nil
}
