package worker

import (
	"github.com/frizz925/higuchi/internal/filter"
)

type Worker struct {
	filters []filter.Filter
}

func New(filters ...filter.Filter) *Worker {
	return &Worker{filters}
}

func (w *Worker) Handle(ctx *filter.Context) error {
	for _, f := range w.filters {
		if err := f.Do(ctx); err != nil {
			return err
		}
	}
	return nil
}
