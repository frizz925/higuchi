package pool

import (
	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/worker"
)

type FixedPool struct {
	workers chan *worker.Worker
}

func NewFixedPool(factory Factory, size int) *FixedPool {
	workers := make(chan *worker.Worker, size)
	for i := 0; i < size; i++ {
		workers <- factory(i)
	}
	return &FixedPool{workers}
}

func (p *FixedPool) Dispatch(ctx *filter.Context, callback Callback) {
	w := <-p.workers
	go func(w *worker.Worker, cb Callback) {
		cb(ctx, w.Handle(ctx))
		p.workers <- w
	}(w, callback)
}
