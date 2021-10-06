package pool

import (
	"sync"

	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/worker"
	"go.uber.org/atomic"
)

type SinglePool struct {
	worker  *worker.Worker
	taskCh  chan singleTask
	stopCh  chan struct{}
	wg      sync.WaitGroup
	running atomic.Bool
}

type singleTask struct {
	ctx      *filter.Context
	callback Callback
}

func NewSinglePool(w *worker.Worker) *SinglePool {
	return &SinglePool{
		worker: w,
		taskCh: make(chan singleTask),
		stopCh: make(chan struct{}, 1),
	}
}

func (p *SinglePool) Start() {
	if p.running.Load() {
		return
	}
	p.running.Store(true)
	p.wg.Add(1)

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

func (p *SinglePool) Stop() {
	if !p.running.Load() {
		return
	}
	p.stopCh <- struct{}{}
	p.wg.Wait()
}

func (p *SinglePool) Dispatch(ctx *filter.Context, callback Callback) {
	p.taskCh <- singleTask{ctx, callback}
}
