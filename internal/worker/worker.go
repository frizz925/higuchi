package worker

import (
	"github.com/frizz925/higuchi/internal/filter"
	"go.uber.org/zap"
)

type Worker struct {
	num     int
	filters []filter.Filter
}

func New(num int, filters ...filter.Filter) *Worker {
	return &Worker{num, filters}
}

func (w *Worker) Handle(ctx *filter.Context) error {
	ctx.Logger = ctx.Logger.With(zap.Int("worker", w.num))
	for _, f := range w.filters {
		if err := f.Do(ctx); err != nil {
			return err
		}
	}
	return nil
}
