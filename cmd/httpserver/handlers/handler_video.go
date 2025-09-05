package handlers

import (
	"log"
	"os"

	"github.com/itsjoeoui/httpfromtcp/internal/headers"
	"github.com/itsjoeoui/httpfromtcp/internal/request"
	"github.com/itsjoeoui/httpfromtcp/internal/response"
)

func HandlerVideo(w *response.Writer, _ *request.Request) {
	err := w.WriteStatusLine(response.StatusCodeOK)
	if err != nil {
		log.Printf("Failed to write status line: %v", err)
	}

	f, err := os.ReadFile("./assets/vim.mp4")
	if err != nil {
		log.Printf("Failed to read file: %v, you can download it with 'just setup'", err)
	}

	h := response.GetDefaultHeaders(len(f))
	h.Override(headers.ContentTypeHeader, "video/mp4")
	err = w.WriteHeaders(h)
	if err != nil {
		log.Printf("Failed to write headers: %v", err)
	}

	_, err = w.WriteBody(f)
	if err != nil {
		log.Printf("Failed to write body: %v", err)
	}
}
