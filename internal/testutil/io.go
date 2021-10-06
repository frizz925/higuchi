package testutil

import "io"

func EchoReadWriter(rw io.ReadWriter) error {
	buf := make([]byte, 65535)
	n, err := rw.Read(buf)
	if err != nil {
		return err
	}
	_, err = rw.Write(buf[:n])
	return err
}
