// Package request implements HTTP request parsing.
package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

var (
	supportedHTTPMethods  = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"}
	supportedHTTPVersions = []string{"1.1"}
)

const crlf = "\r\n"

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HTTPVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestLine, err := parseRequestLine(rawBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request line: %w", err)
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(req []byte) (*RequestLine, error) {
	requestLine := strings.Split(string(req), crlf)[0]
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, errors.New("invalid request line: incorrect number of parts")
	}

	if !slices.Contains(supportedHTTPMethods, parts[0]) {
		return nil, errors.New("invalid request line: unsupported HTTP method")
	}

	httpVersion := strings.TrimPrefix(parts[2], "HTTP/")
	if !slices.Contains(supportedHTTPVersions, httpVersion) {
		return nil, errors.New("invalid request line: unsupported HTTP version")
	}

	return &RequestLine{
		HTTPVersion:   httpVersion,
		RequestTarget: parts[1],
		Method:        parts[0],
	}, nil
}
