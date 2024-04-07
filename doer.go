package worker

import "context"

type Doer interface {
	Do(ctx context.Context) error
}
