package worker

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNoDoers     = errors.New("no doers provided")
	ErrNoExecuters = errors.New("no executers provided")
)

type worker struct {
	doers     []doer
	executers []executer
}

func New(options ...optFunc) (*worker, error) {
	w := worker{}

	for _, opt := range options {
		opt(&w)
	}

	if w.doers == nil {
		return nil, ErrNoDoers
	}

	if w.executers == nil {
		return nil, ErrNoExecuters
	}

	return &w, nil
}

func (w worker) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	doerChannel := make(chan doer, len(w.doers))

	for _, e := range w.executers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for doer := range doerChannel {
				select {
				case <-ctx.Done():
					return
				default:
					_ = e.Execute(ctx, doer.Do)
				}
			}
		}()

	}

	for _, doer := range w.doers {
		doerChannel <- doer
	}
	close(doerChannel)

	wg.Wait()

	if err := ctx.Err(); err != nil {
		return err
	}

	return nil
}
