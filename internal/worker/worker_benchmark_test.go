package worker

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/frizz925/higuchi/internal/crypto/hasher"
	"github.com/frizz925/higuchi/internal/dispatcher"
	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/httputil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const TestBufferSize = 1024

func BenchmarkWorker(b *testing.B) {
	require := require.New(b)
	logger := zap.NewNop()
	username, password := "user", "password"
	magicaddr := "memory:80"

	sampleBody := make([]byte, TestBufferSize)
	n, err := rand.Read(sampleBody)
	require.NoError(err)
	sampleBody = sampleBody[:n]

	h := hasher.NewMD5Hasher([]byte("pepper"))
	md, err := hasher.NewMD5Digest(h, password)
	require.NoError(err)
	users := map[string]interface{}{
		username: md,
	}

	bw := bufio.NewWriterSize(nil, TestBufferSize)
	df := dispatcher.DispatcherFunc(func(rw io.ReadWriter, _ string) error {
		bw.Reset(rw)
		httputil.WriteResponseHeader(&http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
		}, bw)
		return bw.Flush()
	})
	w := New(0,
		filter.NewHealthCheckFilter("OPTIONS", "/healthz"),
		filter.NewParseFilter(TestBufferSize,
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
			filter.NewTunnelFilter(TestBufferSize),
			filter.NewForwardFilter(filter.NewDispatchFilter(df)),
		),
	)

	connCh := make(chan net.Conn)
	defer close(connCh)

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   magicaddr,
				User:   url.UserPassword(username, password),
			}),
			Dial: func(network, addr string) (net.Conn, error) {
				if network != "tcp" || addr != magicaddr {
					return nil, fmt.Errorf("unexpected dial: %s %s", network, addr)
				}
				c1, c2 := net.Pipe()
				connCh <- c2
				return c1, nil
			},
		},
	}
	go func() {
		for {
			conn, ok := <-connCh
			if !ok {
				return
			}
			if err := w.Handle(filter.NewContext(conn, logger)); err != nil {
				b.Logf("error handling: %v", err)
			}
			conn.Close()
		}
	}()

	for i := 0; i < b.N; i++ {
		_, err := client.Post("http://pipe/", "application/octet-stream", bytes.NewReader(sampleBody))
		require.NoError(err)
	}
}
