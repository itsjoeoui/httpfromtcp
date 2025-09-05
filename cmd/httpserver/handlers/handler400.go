package handlers

import (
	"log"
	"os"

	"github.com/itsjoeoui/httpfromtcp/internal/headers"
	"github.com/itsjoeoui/httpfromtcp/internal/request"
	"github.com/itsjoeoui/httpfromtcp/internal/response"
)

func Handler400(w *response.Writer, _ *request.Request) {
	err := w.WriteStatusLine(response.StatusCodeBadRequest)
	if err != nil {
		log.Printf("Failed to write status line: %v", err)
	}

	f, err := os.ReadFile("./cmd/httpserver/templates/400.html")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
	}

	h := response.GetDefaultHeaders(len(f))
	h.Override(headers.ContentTypeHeader, "text/html")
	err = w.WriteHeaders(h)
	if err != nil {
		log.Printf("Failed to write headers: %v", err)
	}

	_, err = w.WriteBody(f)
	if err != nil {
		log.Printf("Failed to write body: %v", err)
	}
}
