package worker_test

import (
	"context"
	"testing"

	"github.com/diego-augusto/go-worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockDoer struct {
	DoMock func(ctx context.Context) error
}

func (m MockDoer) Do(ctx context.Context) error {
	return m.DoMock(ctx)
}

type MockExecuter struct {
	DoExecute func(ctx context.Context, fn func(ctx context.Context) error) error
}

func (m MockExecuter) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.DoExecute(ctx, fn)
}

func TestWorker(t *testing.T) {

	j1 := MockDoer{
		DoMock: func(ctx context.Context) error {
			return nil
		},
	}
	j2 := MockDoer{
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
		worker.WithDoers(j1, j2),
		worker.WithExecuters(e1, e2),
	)
	require.NoError(t, err)

	err = worker.Run(context.Background())
	require.NoError(t, err)

	assert.NotNil(t, worker)
}

func TestWorker_DefaultExecuter(t *testing.T) {

	j1 := MockDoer{
		DoMock: func(ctx context.Context) error {
			return nil
		},
	}
	j2 := MockDoer{
		DoMock: func(ctx context.Context) error {
			return nil
		},
	}

	worker, err := worker.New(
		worker.WithDoers(j1, j2),
		worker.WithDefaultExecuter(),
	)
	require.NoError(t, err)

	err = worker.Run(context.Background())
	require.NoError(t, err)

	assert.NotNil(t, worker)
}

func TestWorker_NoDoers(t *testing.T) {

	w, err := worker.New()
	require.Nil(t, w)
	assert.Equal(t, worker.ErrNoDoers, err)

	w, err = worker.New(worker.WithDoers())
	require.Nil(t, w)
	assert.Equal(t, worker.ErrNoDoers, err)
}

func TestWorker_NoExecuters(t *testing.T) {

	w, err := worker.New(
		worker.WithDoers(MockDoer{}),
		worker.WithExecuters(),
	)
	require.Nil(t, w)
	assert.Equal(t, worker.ErrNoExecuters, err)

	w, err = worker.New(
		worker.WithDoers(MockDoer{}),
		worker.WithExecuters(),
	)
	require.Nil(t, w)
	assert.Equal(t, worker.ErrNoExecuters, err)
}

func TestRunContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	w, err := worker.New(
		worker.WithDoers(MockDoer{}),
		worker.WithDefaultExecuter(),
	)
	require.NoError(t, err)

	err = w.Run(ctx)
	assert.Equal(t, context.Canceled, err)
}
