package pool

import (
	"sync"

	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/worker"
	"go.uber.org/atomic"
)

type DynamicPool struct {
	pool    sync.Pool
	counter atomic.Int64
}

func NewDynamicPool(f Factory) *DynamicPool {
	dp := &DynamicPool{}
	dp.pool.New = func() interface{} {
		w := f(int(dp.counter.Load()))
		dp.counter.Add(1)
		return w
	}
	return dp
}

func (p *DynamicPool) Dispatch(ctx *filter.Context, callback Callback) {
	w := p.pool.Get().(*worker.Worker)
	go p.dispatch(w, ctx, callback)
}

func (p *DynamicPool) dispatch(w *worker.Worker, ctx *filter.Context, callback Callback) {
	err := w.Handle(ctx)
	p.pool.Put(w)
	callback(ctx, err)
}
