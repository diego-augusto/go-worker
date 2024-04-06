package worker_test

import (
	"context"
	"testing"

	"github.com/diego-augusto/go-worker"
)

type MockExecuter struct {
	DoMock func(ctx context.Context) error
}

func (m MockExecuter) Do(ctx context.Context) error {
	return m.DoMock(ctx)
}

func TestWorker(t *testing.T) {

	e1 := MockExecuter{
		DoMock: func(ctx context.Context) error {
			return nil
		},
	}
	e2 := MockExecuter{
		DoMock: func(ctx context.Context) error {
			return nil
		},
	}

	worker := worker.NewWorker(e1, e2)

	err := worker.Run(context.Background(), 1)
	if err != nil {
		t.Fatalf("Error running worker: %v", err)
	}
}
