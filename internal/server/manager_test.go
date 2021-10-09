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

type managerTestSuite struct {
	suite.Suite

	client  *http.Client
	pool    *pool.SinglePool
	manager *Manager
	server  *Server
}

func TestManager(t *testing.T) {
	suite.Run(t, new(managerTestSuite))
}

func (ts *managerTestSuite) SetupSuite() {
	var err error
	require := ts.Require()
	logger := zap.NewExample()
	user, pass := "testuser", "testpass"

	ts.pool = pool.NewSinglePool(worker.New(0,
		filter.NewParseFilter(
			DefaultBufferSize,
			filter.NewCertbotFilter(filter.CertbotConfig{ChallengePath: "/no-challenge"}),
			filter.NewAuthFilter(map[string]interface{}{user: pass}),
			filter.NewTunnelFilter(DefaultBufferSize),
			filter.NewForwardFilter(filter.NewDispatchFilter(dispatcher.NewTCPDispatcher(DefaultBufferSize))),
		),
	))
	ts.pool.Start()

	ts.manager = NewManager(ManagerConfig{
		Pool:   ts.pool,
		Logger: logger,
	})

	ts.server, err = ts.manager.ListenAndServe("tcp", "127.0.0.1:0")
	require.NoError(err)
	proxyUrl := &url.URL{
		Host: ts.server.Addr().String(),
		User: url.UserPassword(user, pass),
	}

	ts.client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
}

func (ts *managerTestSuite) TearDownSuite() {
	require := ts.Require()
	require.NoError(ts.server.Close())
	require.NoError(ts.manager.Close())
	ts.pool.Stop()
}

func (ts *managerTestSuite) TestE2E() {
	require := ts.Require()
	srv, expected, err := createTestWebServer()
	require.NoError(err)
	defer srv.Close()
	res, err := ts.client.Get(srv.URL)
	require.NoError(err)
	defer res.Body.Close()
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
