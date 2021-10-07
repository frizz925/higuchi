package pool

import (
	"testing"

	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/testutil"
	"github.com/frizz925/higuchi/internal/worker"
)

func TestPreallocatedPool(t *testing.T) {
	p := NewPreallocatedPool(func(num int) *worker.Worker {
		return worker.New(num, filter.FilterFunc(func(c *filter.Context, _ filter.Next) error {
			return testutil.EchoReadWriter(c)
		}))
	}, 1)
	testPool(t, p)
}
