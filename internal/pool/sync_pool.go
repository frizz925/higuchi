package pool

import (
	"net"

	"github.com/frizz925/higuchi/internal/worker"
)

type SyncPool struct {
	worker *worker.Worker
}

func NewSyncPool(w *worker.Worker) *SyncPool {
	return &SyncPool{w}
}

func (p *SyncPool) Dispatch(conn net.Conn, callback Callback) {
	callback(conn, p.worker.Handle(conn))
}
