package filter

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestForwardedFilter(t *testing.T) {
	require := require.New(t)
	ctx := NewContext(nil, zap.NewExample())
	err := DefaultForwardedFilter.Do(ctx, &http.Request{
		Header: http.Header{
			"X-Forwarded-Proto": []string{"dummy"},
			"X-Forwarded-For":   []string{"localhost:20450"},
			"X-Forwarded-Host":  []string{"localhost:8080"},
		},
	}, NextNoop)
	require.NoError(err)
	require.Equal("dummy", ctx.LogFields.Proto)
	require.Equal("localhost:20450", ctx.LogFields.Source)
	require.Equal("localhost:8080", ctx.LogFields.Server)
}
