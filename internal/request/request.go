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

const (
	crlf       = "\r\n"
	bufferSize = 8
)

type Request struct {
	RequestLine RequestLine
	ParserState ParserState
}

type ParserState string

const (
	ParserStateInitialized ParserState = "Initialized"
	ParserStateDone        ParserState = "Done"
)

type RequestLine struct {
	HTTPVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.ParserState {
	case ParserStateInitialized:
		requestLine, length, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if length == 0 {
			// this means we need more data to parse the request line
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.ParserState = ParserStateDone
		return length, nil
	case ParserStateDone:
		return 0, fmt.Errorf("request already fully parsed")
	default:
		return 0, fmt.Errorf("unknown parser state: %s", r.ParserState)
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		ParserState: ParserStateInitialized,
	}

	buffer := make([]byte, bufferSize)
	readToIndex := 0

	for request.ParserState != ParserStateDone {
		if readToIndex == len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		bytesRead, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if bytesRead != 0 {
					panic("What the hell, I thought if io.EOF is returned, bytesRead should be 0")
				}
				request.ParserState = ParserStateDone
				break
			}
			return nil, err
		}

		readToIndex += bytesRead

		parsedToIndex, err := request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		if parsedToIndex != 0 {
			copy(buffer, buffer[parsedToIndex:])
			readToIndex -= parsedToIndex
		}
	}

	return request, nil
}

func parseRequestLine(req []byte) (*RequestLine, int, error) {
	splitReq := strings.Split(string(req), crlf)
	if len(splitReq) <= 1 {
		// we do not have a complete request line yet
		return nil, 0, nil
	}

	requestLine := splitReq[0]

	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, 0, errors.New("invalid request line: incorrect number of parts")
	}

	if !slices.Contains(supportedHTTPMethods, parts[0]) {
		return nil, 0, errors.New("invalid request line: unsupported HTTP method")
	}

	httpVersion := strings.TrimPrefix(parts[2], "HTTP/")
	if !slices.Contains(supportedHTTPVersions, httpVersion) {
		return nil, 0, errors.New("invalid request line: unsupported HTTP version")
	}

	return &RequestLine{
		HTTPVersion:   httpVersion,
		RequestTarget: parts[1],
		Method:        parts[0],
	}, len(requestLine) + len(crlf), nil
}
