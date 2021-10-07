package httputil

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

func WriteResponseHeader(res *http.Response, bw *bufio.Writer) {
	proto := res.Proto
	if proto == "" {
		proto = fmt.Sprintf("HTTP/%d.%d", res.ProtoMajor, res.ProtoMinor)
	}
	bw.WriteString(proto + " ")

	status := res.Status
	if status == "" {
		code := res.StatusCode
		if code <= 0 {
			code = http.StatusOK
		}
		status = fmt.Sprintf("%d %s", code, http.StatusText(code))
	}
	bw.WriteString(status + "\r\n")

	for k, v := range res.Header {
		bw.WriteString(fmt.Sprintf("%s: %s\r\n", k, strings.Join(v, ", ")))
	}
	bw.WriteString("\r\n")
}
