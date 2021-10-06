package ioutil

import "io"

type PrefixedReader struct {
	io.Reader
	b []byte
}

func NewPrefixedReader(b []byte, r io.Reader) *PrefixedReader {
	return &PrefixedReader{r, b}
}

func (pr *PrefixedReader) Read(b []byte) (int, error) {
	if len(pr.b) <= 0 {
		return pr.Reader.Read(b)
	}
	n := copy(b, pr.b)
	pr.b = pr.b[n:]
	return n, nil
}
