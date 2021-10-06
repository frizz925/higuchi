package pool

import (
	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/worker"
)

type PreallocatedPool struct {
	workers chan *worker.Worker
}

func NewPreallocatedPool(factory Factory, size int) *PreallocatedPool {
	if size < 1 {
		size = 1
	}
	workers := make(chan *worker.Worker, size)
	for i := 0; i < size; i++ {
		workers <- factory(i)
	}
	return &PreallocatedPool{workers}
}

func (p *PreallocatedPool) Dispatch(ctx *filter.Context, callback Callback) {
	w := <-p.workers
	go func(w *worker.Worker, cb Callback) {
		cb(ctx, w.Handle(ctx))
		p.workers <- w
	}(w, callback)
}
