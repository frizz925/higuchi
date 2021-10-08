package filter

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAuthFilter(t *testing.T) {
	require := require.New(t)
	user := "testuser"

	buf := make([]byte, 8)
	n, err := rand.Read(buf)
	require.NoError(err)
	pass := hex.EncodeToString(buf[:n])

	authParam := fmt.Sprintf("%s:%s", user, pass)
	authParam = base64.StdEncoding.EncodeToString([]byte(authParam))
	header := make(http.Header)

	af := NewAuthFilter(map[string]interface{}{
		user: pass,
	})

	conn, _ := net.Pipe()
	ctx := NewContext(conn, zap.NewExample())
	require.Error(af.Do(ctx, &http.Request{Header: header}, NextNoop))
	header.Set("Proxy-Authorization", "Basic ")
	require.Error(af.Do(ctx, &http.Request{Header: header}, NextNoop))
	header.Set("Proxy-Authorization", "Basic userpass")
	require.Error(af.Do(ctx, &http.Request{Header: header}, NextNoop))
	header.Set("Proxy-Authorization", "Basic aGVsbG8K")
	require.Error(af.Do(ctx, &http.Request{Header: header}, NextNoop))
	header.Set("Proxy-Authorization", "Basic "+authParam)
	require.NoError(af.Do(ctx, &http.Request{Header: header}, NextNoop))
	require.Equal(user, ctx.LogFields.User)
}
