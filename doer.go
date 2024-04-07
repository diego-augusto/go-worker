package worker

import "context"

type doer interface {
	Do(ctx context.Context) error
}
