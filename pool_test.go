package worker_test

import (
	"context"
	"testing"

	"github.com/diego-augusto/go-worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockWorker struct {
	DoMock func(ctx context.Context) error
}

func (m MockWorker) Do(ctx context.Context) error {
	return m.DoMock(ctx)
}

type MockExecuter struct {
	DoExecute func(ctx context.Context, fn func(ctx context.Context) error) error
}

func (m MockExecuter) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.DoExecute(ctx, fn)
}

func TestWorker(t *testing.T) {

	j1 := MockWorker{
		DoMock: func(ctx context.Context) error {
			return nil
		},
	}
	j2 := MockWorker{
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

	worker, err := worker.NewPool(
		worker.WithWorkers(j1, j2),
		worker.WithExecuters(e1, e2),
	)
	require.NoError(t, err)

	err = worker.Run(context.Background())
	require.NoError(t, err)

	assert.NotNil(t, worker)
}

func TestWorker_DefaultExecuter(t *testing.T) {

	j1 := MockWorker{
		DoMock: func(ctx context.Context) error {
			return nil
		},
	}
	j2 := MockWorker{
		DoMock: func(ctx context.Context) error {
			return nil
		},
	}

	worker, err := worker.NewPool(
		worker.WithWorkers(j1, j2),
		worker.WithDefaultExecuter(),
	)
	require.NoError(t, err)

	err = worker.Run(context.Background())
	require.NoError(t, err)

	assert.NotNil(t, worker)
}

func TestWorker_NoWorkers(t *testing.T) {

	w, err := worker.NewPool()
	require.Nil(t, w)
	assert.Equal(t, worker.ErrNoWorkers, err)

	w, err = worker.NewPool(worker.WithWorkers())
	require.Nil(t, w)
	assert.Equal(t, worker.ErrNoWorkers, err)
}

func TestWorker_NoExecuters(t *testing.T) {

	w, err := worker.NewPool(
		worker.WithWorkers(MockWorker{}),
		worker.WithExecuters(),
	)
	require.Nil(t, w)
	assert.Equal(t, worker.ErrNoExecuters, err)

	w, err = worker.NewPool(
		worker.WithWorkers(MockWorker{}),
		worker.WithExecuters(),
	)
	require.Nil(t, w)
	assert.Equal(t, worker.ErrNoExecuters, err)
}

func TestRunContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	w, err := worker.NewPool(
		worker.WithWorkers(MockWorker{}),
		worker.WithDefaultExecuter(),
	)
	require.NoError(t, err)

	err = w.Run(ctx)
	assert.Equal(t, context.Canceled, err)
}
