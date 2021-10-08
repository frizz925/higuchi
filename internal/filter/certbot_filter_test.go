package filter

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type certbotFilterTestSuite struct {
	suite.Suite
	hostname      string
	webroot       string
	challengePath string
	filename      string
	expected      []byte
}

func TestCertbotFilter(t *testing.T) {
	suite.Run(t, &certbotFilterTestSuite{
		hostname:      "certbot-test",
		webroot:       "/tmp",
		challengePath: "/certbot-filter-test",
		filename:      "certbot-test.txt",
		expected:      []byte("expected"),
	})
}

func (s *certbotFilterTestSuite) SetupTest() {
	require := s.Require()
	dir := path.Join(s.webroot, s.challengePath)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		require.NoError(os.MkdirAll(dir, 0755))
	}
	fn := path.Join(dir, s.filename)
	require.NoError(os.WriteFile(fn, s.expected, 0644))
}

func (s *certbotFilterTestSuite) TestGet() {
	require := s.Require()
	cf := NewCertbotFilter(CertbotConfig{
		Hostname:      s.hostname,
		Webroot:       s.webroot,
		ChallengePath: s.challengePath,
	})
	host := fmt.Sprintf("%s:80", s.hostname)
	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "http",
			Host:   host,
			Path:   path.Join(s.challengePath, s.filename),
		},
		Host: host,
	}

	c1, c2 := net.Pipe()
	defer c1.Close()

	errCh := make(chan error, 1)
	go func(ctx *Context, req *http.Request) {
		defer ctx.Close()
		errCh <- cf.Do(ctx, req, func() error {
			return errors.New("unexpected next")
		})
		close(errCh)
	}(NewContext(c2, zap.NewExample()), req)

	res, err := http.ReadResponse(bufio.NewReader(c1), req)
	require.NoError(<-errCh)
	require.NoError(err)
	require.Equal(http.StatusOK, res.StatusCode)
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(err)
	require.Equal(s.expected, b)
}
