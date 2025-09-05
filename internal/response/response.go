// Package response provides utilities for writing HTTP responses.
package response

import (
	"fmt"
	"io"

	"github.com/itsjoeoui/httpfromtcp/internal/common"
	"github.com/itsjoeoui/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusCodeOK                  StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	reasonPhrase, ok := statusCodeToReasonPhrase[statusCode]
	if !ok {
		reasonPhrase = "" // just leave it blank if unknown
	}

	_, err := fmt.Fprintf(w.writer, "HTTP/1.1 %d %s%s", statusCode, reasonPhrase, common.CRLF)
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for k, v := range headers {
		_, err := fmt.Fprintf(w.writer, "%s: %s%s", k, v, common.CRLF)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(w.writer, common.CRLF)
	return err
}

func (w *Writer) WriteBody(body []byte) (int, error) {
	bytesWritten, err := fmt.Fprintf(w.writer, "%s", body)
	if err != nil {
		return 0, err
	}

	return bytesWritten, nil
}

var statusCodeToReasonPhrase map[StatusCode]string = map[StatusCode]string{
	StatusCodeOK:                  "OK",
	StatusCodeBadRequest:          "Bad Request",
	StatusCodeInternalServerError: "Internal Server Error",
}

func GetDefaultHeaders(contentLength int) headers.Headers {
	h := headers.Headers{}

	h.Set(headers.ContentTypeHeader, "text/plain")
	h.Set(headers.ContentLengthHeader, fmt.Sprintf("%d", contentLength))
	h.Set(headers.ConnectionHeader, "close")

	return h
}
