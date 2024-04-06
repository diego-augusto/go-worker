package worker

import "context"

type Executer interface {
	Execute(ctx context.Context, fn func(ctx context.Context) error) error
}

type defaultExecuter struct {
}

func NewDefaultExecuter() Executer {
	return defaultExecuter{}
}

func (e defaultExecuter) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
