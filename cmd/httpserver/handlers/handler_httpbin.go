package handlers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/itsjoeoui/httpfromtcp/internal/headers"
	"github.com/itsjoeoui/httpfromtcp/internal/request"
	"github.com/itsjoeoui/httpfromtcp/internal/response"
)

func HandlerHTTPBin(w *response.Writer, r *request.Request) {
	route := strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin")
	resp, err := http.Get(fmt.Sprintf("https://httpbin.org%s", route))
	if err != nil {
		log.Printf("Failed to fetch from httpbin: %v", err)
		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("Failed to close httpbin response body: %v", err)
		}
	}()

	err = w.WriteStatusLine(response.StatusCodeOK)
	if err != nil {
		log.Printf("Failed to write status line: %v", err)
	}

	h := response.GetDefaultHeaders(0)
	h.Remove(headers.ContentLengthHeader)
	h.Override(headers.TransferEncodingHeader, "chunked")
	h.Set(headers.TrailerHeader, headers.XContentLengthHeader)
	h.Set(headers.TrailerHeader, headers.XContentSHA256)

	err = w.WriteHeaders(h)
	if err != nil {
		log.Printf("Failed to write headers: %v", err)
	}

	fullBody := make([]byte, 0)

	buffer := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			_, writeErr := w.WriteChunkedBody(buffer[:n])
			if writeErr != nil {
				log.Printf("Failed to write chunked body: %v", writeErr)
			}

			fullBody = append(fullBody, buffer[:n]...)
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read from httpbin response body: %v", err)
			}
			break
		}
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Printf("Failed to write chunked body done: %v", err)
	}

	trailers := headers.NewHeaders()
	sha256 := fmt.Sprintf("%x", sha256.Sum256(fullBody))
	trailers.Override(headers.XContentSHA256, sha256)
	trailers.Override(headers.XContentLengthHeader, fmt.Sprintf("%d", len(fullBody)))
	err = w.WriteTrailers(trailers)
	if err != nil {
		log.Printf("Failed to write trailers: %v", err)
	}
}
