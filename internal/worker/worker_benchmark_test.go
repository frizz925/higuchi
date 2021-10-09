package worker

import (
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/frizz925/higuchi/internal/crypto/hasher"
	"github.com/frizz925/higuchi/internal/dispatcher"
	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/testutil"

	"go.uber.org/zap"
)

const testBufferSize = 512

var (
	rawProxyRequest = testutil.LinesToRawHeader(
		"GET http://pipe/ HTTP/1.1",
		"Host: pipe",
		"Proxy-Authorization: Basic dXNlcjpwYXNzd29yZA==",
	)
	rawProxyResponse = testutil.LinesToRawHeader(
		"HTTP/1.1 200 OK",
	)
)

func BenchmarkWorker(b *testing.B) {
	logger := zap.NewNop()
	username, password := "user", "password"

	h := hasher.NewMD5Hasher([]byte("pepper"))
	md, err := hasher.NewMD5Digest(h, password)
	if err != nil {
		b.Fatal("error while creating digest", err)
	}
	users := map[string]interface{}{
		username: md,
	}

	df := dispatcher.DispatcherFunc(func(rw io.ReadWriter, _ string) error {
		_, err := rw.Write(rawProxyResponse)
		return err
	})

	filters := []filter.Filter{
		filter.NewHealthCheckFilter("OPTIONS", "/healthz"),
		filter.NewParseFilter(testBufferSize,
			filter.DefaultForwardedFilter,
			filter.NewCertbotFilter(filter.CertbotConfig{
				Hostname:      "localhost",
				Webroot:       "/dev/null",
				ChallengePath: "/no-challenge",
			}),
			filter.NewAuthFilter(users, func(password string, i interface{}) bool {
				v, ok := i.(hasher.PasswordDigest)
				if !ok {
					return false
				}
				return v.Compare(password) == 0
			}),
			filter.NewTunnelFilter(testBufferSize),
			filter.NewForwardFilter(filter.NewDispatchFilter(df)),
		),
	}

	var counter uint32 = 0
	pool := sync.Pool{
		New: func() interface{} {
			num := atomic.AddUint32(&counter, 1)
			return New(int(num), filters...)
		},
	}

	connCh := make(chan net.Conn)
	doneCh := make(chan struct{})
	go func(connCh <-chan net.Conn, doneCh <-chan struct{}) {
		for {
			select {
			case conn := <-connCh:
				w := pool.Get().(*Worker)
				ctx := filter.NewContext(conn, logger)
				if err := w.Handle(ctx); err != nil {
					b.Logf("error handling: %v", err)
				}
				conn.Close()
				pool.Put(w)
			case <-doneCh:
				return
			}
		}
	}(connCh, doneCh)

	buf := make([]byte, testBufferSize)
	for i := 0; i < b.N; i++ {
		c1, c2 := net.Pipe()
		connCh <- c2
		_, err = c1.Write(rawProxyRequest)
		if err != nil {
			b.Fatal("error while writing request", err)
		}
		_, err = c1.Read(buf)
		if err != nil {
			b.Fatal("error while reading response", err)
		}
	}
}
