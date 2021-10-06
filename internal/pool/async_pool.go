package pool

import (
	"sync"

	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/worker"
	"go.uber.org/atomic"
)

type AsyncPool struct {
	worker  *worker.Worker
	taskCh  chan asyncTask
	stopCh  chan struct{}
	wg      sync.WaitGroup
	running atomic.Bool
}

type asyncTask struct {
	ctx      *filter.Context
	callback Callback
}

func NewAsyncPool(w *worker.Worker) *AsyncPool {
	return &AsyncPool{
		worker: w,
		taskCh: make(chan asyncTask),
		stopCh: make(chan struct{}, 1),
	}
}

func (p *AsyncPool) Start() {
	if p.running.Load() {
		return
	}
	p.running.Store(true)

	go func() {
		defer p.wg.Done()
		for {
			select {
			case task := <-p.taskCh:
				task.callback(task.ctx, p.worker.Handle(task.ctx))
			case <-p.stopCh:
				return
			}
		}
	}()
}

func (p *AsyncPool) Stop() {
	if !p.running.Load() {
		return
	}
	p.stopCh <- struct{}{}
	p.wg.Wait()
}

func (p *AsyncPool) Dispatch(ctx *filter.Context, callback Callback) {
	p.taskCh <- asyncTask{ctx, callback}
}
