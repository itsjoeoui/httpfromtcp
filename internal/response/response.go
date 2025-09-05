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
	Writer io.Writer
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	reasonPhrase, ok := statusCodeToReasonPhrase[statusCode]
	if !ok {
		reasonPhrase = "" // just leave it blank if unknown
	}

	_, err := fmt.Fprintf(w.Writer, "HTTP/1.1 %d %s%s", statusCode, reasonPhrase, common.CRLF)
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for k, v := range headers {
		_, err := fmt.Fprintf(w.Writer, "%s: %s%s", k, v, common.CRLF)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(w.Writer, common.CRLF)
	return err
}

func (w *Writer) WriteBody(body []byte) (int, error) {
	bytesWritten, err := fmt.Fprintf(w.Writer, "%s", body)
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

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reasonPhrase, ok := statusCodeToReasonPhrase[statusCode]
	if !ok {
		reasonPhrase = "" // just leave it blank if unknown
	}

	_, err := fmt.Fprintf(w, "HTTP/1.1 %d %s%s", statusCode, reasonPhrase, common.CRLF)
	return err
}

func GetDefaultHeaders(contentLength int) headers.Headers {
	h := headers.Headers{}

	h.Set(headers.ContentTypeHeader, "text/plain")
	h.Set(headers.ContentLengthHeader, fmt.Sprintf("%d", contentLength))
	h.Set(headers.ConnectionHeader, "close")

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := fmt.Fprintf(w, "%s: %s%s", k, v, common.CRLF)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(w, common.CRLF)
	return err
}
