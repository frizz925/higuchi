package dispatcher

import (
	"io"
	"testing"
)

func TestDispatcherFunc(t *testing.T) {
	d := DispatcherFunc(func(rw io.ReadWriter, addr string) error {
		return nil
	})
	if err := d.Dispatch(nil, ""); err != nil {
		t.Fatal(err)
	}
}
