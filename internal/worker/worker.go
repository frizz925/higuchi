package worker

import (
	"github.com/frizz925/higuchi/internal/filter"
)

type Worker struct {
	num     int
	filters []filter.Filter
}

func New(num int, filters ...filter.Filter) *Worker {
	return &Worker{num, filters}
}

func (w *Worker) Handle(ctx *filter.Context) error {
	ctx.LogFields.Worker = w.num
	ctx.UpdateLogger()

	var next filter.Next
	idx := 0
	next = func() error {
		if idx >= len(w.filters) {
			return nil
		}
		f := w.filters[idx]
		idx++
		return f.Do(ctx, next)
	}
	return next()
}
