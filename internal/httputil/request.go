package httputil

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
)

func ParseRequestAddress(req *http.Request) string {
	hostport := req.Host
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		host = hostport
		port = "80"
	}
	return net.JoinHostPort(host, port)
}

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
