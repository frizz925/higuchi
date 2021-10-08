package pool

import (
	"testing"

	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/testutil"
	"github.com/frizz925/higuchi/internal/worker"
)

func TestDynamicPool(t *testing.T) {
	p := NewDynamicPool(func(num int) *worker.Worker {
		return worker.New(num, filter.FilterFunc(func(c *filter.Context, _ filter.Next) error {
			return testutil.EchoReadWriter(c)
		}))
	})
	testPool(t, p)
}
