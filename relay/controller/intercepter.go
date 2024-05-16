package controller

import (
	"bytes"
	"io"
	"net/http"
)

type AuditReadCloser struct {
	Reader io.Reader
	Closer io.Closer
	Buffer *bytes.Buffer
}

func (arc *AuditReadCloser) Read(p []byte) (int, error) {
	n, err := arc.Reader.Read(p)
	if n > 0 {
		arc.Buffer.Write(p[:n])
	}
	return n, err
}

func (arc *AuditReadCloser) Close() error {
	return arc.Closer.Close()
}

func captureResponseBody(resp *http.Response) *bytes.Buffer {
	buf := &bytes.Buffer{}
	arc := &AuditReadCloser{
		Reader: resp.Body,
		Closer: resp.Body,
		Buffer: buf,
	}
	resp.Body = arc
	return buf
}
