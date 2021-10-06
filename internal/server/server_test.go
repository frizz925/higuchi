package server

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/frizz925/higuchi/internal/dispatcher"
	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/pool"
	"github.com/frizz925/higuchi/internal/worker"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type serverTestSuite struct {
	suite.Suite

	client   *http.Client
	pool     *pool.AsyncPool
	server   *Server
	listener *Listener
}

func TestServer(t *testing.T) {
	suite.Run(t, new(serverTestSuite))
}

func (ts *serverTestSuite) SetupSuite() {
	var err error
	require := ts.Require()
	logger := zap.NewExample()

	ts.pool = pool.NewAsyncPool(worker.New(0,
		filter.NewParseFilter(filter.NewForwardFilter(
			filter.NewDispatchFilter(dispatcher.NewTCPDispatcher(DefaultBufferSize)),
		)),
	))
	ts.pool.Start()

	ts.server = New(Config{
		Pool:   ts.pool,
		Logger: logger,
	})

	ts.listener, err = ts.server.Listen("tcp", "127.0.0.1:0")
	require.NoError(err)
	proxyUrl := &url.URL{
		Host: ts.listener.Addr().String(),
	}

	ts.client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
}

func (ts *serverTestSuite) TearDownSuite() {
	require := ts.Require()
	require.NoError(ts.listener.Close())
	require.NoError(ts.server.Close())
	ts.pool.Stop()
}

func (ts *serverTestSuite) TestE2E() {
	require := ts.Require()
	srv, expected, err := createTestWebServer()
	require.NoError(err)
	res, err := ts.client.Get(srv.URL)
	require.NoError(err)
	resBody, err := ioutil.ReadAll(res.Body)
	require.NoError(err)
	require.Equal(expected, resBody)
}

func createTestWebServer() (*httptest.Server, []byte, error) {
	expected := make([]byte, 32)
	n, err := rand.Read(expected)
	expected = expected[:n]
	if err != nil {
		return nil, nil, err
	}
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/octet-stream")
		_, _ = rw.Write(expected)
	}))
	return srv, expected, nil
}
