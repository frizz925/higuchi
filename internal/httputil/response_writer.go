package httputil

import (
	"bufio"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"
)

type ResponseWriter struct {
	Response *http.Response

	bw            *bufio.Writer
	cw            io.WriteCloser
	wroteHeader   bool
	contentLength int
}

var _ http.ResponseWriter = (*ResponseWriter)(nil)

func NewResponseWriter(w io.Writer, req *http.Request) *ResponseWriter {
	return &ResponseWriter{
		Response: &http.Response{
			Request:    req,
			StatusCode: http.StatusOK,
			Proto:      req.Proto,
			ProtoMajor: req.ProtoMajor,
			ProtoMinor: req.ProtoMinor,
			Header:     make(http.Header),
			TLS:        req.TLS,
		},
		bw: bufio.NewWriter(w),
		cw: httputil.NewChunkedWriter(w),
	}
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.Response.Header
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.Response.StatusCode = code
	if cl := rw.Response.Header.Get("Content-Length"); cl != "" {
		n, err := strconv.Atoi(cl)
		if err == nil {
			rw.contentLength = n
			return
		}
	}
	rw.Response.Header.Set("Transfer-Encoding", "chunked")
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
		WriteResponseHeader(rw.Response, rw.bw)
		rw.wroteHeader = true
	}

	if rw.contentLength <= 0 && len(b) > 0 {
		if err := rw.bw.Flush(); err != nil {
			return 0, err
		}
		return rw.cw.Write(b)
	}

	n, err := rw.bw.Write(b)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (rw *ResponseWriter) Close() error {
	if rw.contentLength <= 0 {
		if err := rw.cw.Close(); err != nil {
			return err
		}
		rw.bw.WriteString("\r\n")
	}
	if err := rw.bw.Flush(); err != nil {
		return err
	}
	return nil
}
