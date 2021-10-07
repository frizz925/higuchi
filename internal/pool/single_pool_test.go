package pool

import (
	"testing"

	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/testutil"
	"github.com/frizz925/higuchi/internal/worker"
)

func TestSinglePool(t *testing.T) {
	p := NewSinglePool(worker.New(0, filter.FilterFunc(func(c *filter.Context, _ filter.Next) error {
		return testutil.EchoReadWriter(c)
	})))
	p.Start()
	defer p.Stop()
	testPool(t, p)
}
