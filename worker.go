package worker

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNoJobs      = errors.New("no jobs provided")
	ErrNoExecuters = errors.New("no executers provided")
)

type optFunc func(*Worker)

type Job interface {
	Do(ctx context.Context) error
}

func WithJobs(jobs []Job) optFunc {
	return func(w *Worker) {
		w.Jobs = jobs
	}
}

func WithExecuters(executers []Executer) optFunc {
	return func(w *Worker) {
		w.Executers = executers
	}
}

type Worker struct {
	Jobs      []Job
	Executers []Executer
}

func New(options ...optFunc) (*Worker, error) {
	worker := Worker{
		Jobs:      nil,
		Executers: []Executer{NewDefaultExecuter()},
	}

	for _, opt := range options {
		opt(&worker)
	}

	if worker.Jobs == nil {
		return nil, ErrNoJobs
	}

	if worker.Executers == nil {
		return nil, ErrNoExecuters
	}

	return &worker, nil
}

func (w Worker) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	jobChannel := make(chan Job, len(w.Jobs))

	for _, e := range w.Executers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobChannel {
				select {
				case <-ctx.Done():
					return
				default:
					_ = e.Execute(ctx, job.Do)
				}
			}
		}()
	}

	for _, job := range w.Jobs {
		jobChannel <- job
	}
	close(jobChannel)

	wg.Wait()

	return nil
}
