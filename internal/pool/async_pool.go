package pool

import (
	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/worker"
	"go.uber.org/atomic"
)

type AsyncPool struct {
	factory Factory
	counter atomic.Uint64
}

func NewAsyncPool(f Factory) *AsyncPool {
	return &AsyncPool{
		factory: f,
	}
}

func (p *AsyncPool) Dispatch(ctx *filter.Context, callback Callback) {
	num := p.counter.Add(1)
	w := p.factory(int(num))
	go func(w *worker.Worker, cb Callback) {
		cb(ctx, w.Handle(ctx))
	}(w, callback)
}
