package httputil

import (
	"bufio"
	"bytes"
	"net/http"
)

func ParseRequestHeader(br *bufio.Reader) (req *http.Request, err error) {
	var (
		buf  bytes.Buffer
		line []byte
	)

	for isReading := true; isReading; {
		isPrefix := true
		for isPrefix {
			line, isPrefix, err = br.ReadLine()
			if err != nil {
				return nil, err
			} else if len(line) <= 0 {
				isReading = false
				break
			}
			buf.Write(line)
		}
		buf.WriteString("\r\n")
	}
	buf.WriteString("\r\n")

	return http.ReadRequest(bufio.NewReader(&buf))
}
