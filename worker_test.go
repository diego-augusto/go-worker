package worker_test

import (
	"context"
	"testing"

	"github.com/diego-augusto/go-worker"
)

type MockJob struct {
	DoMock func(ctx context.Context) error
}

func (m MockJob) Do(ctx context.Context) error {
	return m.DoMock(ctx)
}

type MockExecuter struct {
	DoExecute func(ctx context.Context, fn func(ctx context.Context) error) error
}

func (m MockExecuter) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.DoExecute(ctx, fn)
}

func TestWorker(t *testing.T) {

	j1 := MockJob{
		DoMock: func(ctx context.Context) error {
			return nil
		},
	}
	j2 := MockJob{
		DoMock: func(ctx context.Context) error {
			return nil
		},
	}

	e1 := MockExecuter{
		DoExecute: func(ctx context.Context, fn func(ctx context.Context) error) error {
			return nil
		},
	}
	e2 := MockExecuter{
		DoExecute: func(ctx context.Context, fn func(ctx context.Context) error) error {
			return nil
		},
	}

	worker, err := worker.New(
		worker.WithJobs([]worker.Job{j1, j2}),
		worker.WithExecuters([]worker.Executer{e1, e2}),
	)
	if err != nil {
		t.Fatalf("Error creating worker: %v", err)
	}

	err = worker.Run(context.Background())
	if err != nil {
		t.Fatalf("Error running worker: %v", err)
	}
}
