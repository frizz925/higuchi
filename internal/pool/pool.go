package pool

import (
	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/worker"
)

type Factory func(num int) *worker.Worker

type Callback func(ctx *filter.Context, err error)

type Pool interface {
	Dispatch(ctx *filter.Context, callback Callback)
}
