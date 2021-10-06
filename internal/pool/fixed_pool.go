package pool

import (
	"net"

	"github.com/frizz925/higuchi/internal/worker"
)

type FixedPool struct {
	workers chan *worker.Worker
}

func NewFixedPool(factory Factory, size int) *FixedPool {
	workers := make(chan *worker.Worker, size)
	for i := 0; i < size; i++ {
		workers <- factory()
	}
	return &FixedPool{workers}
}

func (p *FixedPool) Dispatch(conn net.Conn, callback Callback) {
	w := <-p.workers
	go func(w *worker.Worker) {
		callback(conn, w.Handle(conn))
		p.workers <- w
	}(w)
}
