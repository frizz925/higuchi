package pool

import (
	"net"

	"github.com/frizz925/higuchi/internal/worker"
)

type Factory func() *worker.Worker

type Callback func(conn net.Conn, err error)

type Pool interface {
	Dispatch(conn net.Conn, callback Callback)
}
