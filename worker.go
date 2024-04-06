package worker

import (
	"context"
	"sync"
)

type Executer interface {
	Do(ctx context.Context) error
}

type Worker struct {
	Jobs []Executer
}

func NewWorker(jobs ...Executer) Worker {
	return Worker{Jobs: jobs}
}

func (w Worker) Run(ctx context.Context, workers int) error {
	var wg sync.WaitGroup
	jobChannel := make(chan Executer, len(w.Jobs))

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobChannel {
				select {
				case <-ctx.Done():
					return
				default:
					//TODO: return error
					_ = job.Do(ctx)
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
