// Package request implements HTTP request parsing.
package request

import (
	"errors"
	"io"
	"slices"
	"strings"

	"github.com/itsjoeoui/httpfromtcp/internal/headers"
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
	Headers     headers.Headers
	ParserState ParserState
}

type RequestLine struct {
	HTTPVersion   string
	RequestTarget string
	Method        string
}

type ParserState string

const (
	ParserStateRequestLine ParserState = "RequestLine"
	ParserStateHeaders     ParserState = "Headers"
	ParserStateDone        ParserState = "Done"
)

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for r.ParserState != ParserStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.ParserState {
	case ParserStateRequestLine:
		requestLine, length, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if length == 0 {
			// this means we need more data to parse the request line
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.ParserState = ParserStateHeaders
		return length, nil
	case ParserStateHeaders:
		bytesParsed, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.ParserState = ParserStateDone
		}
		return bytesParsed, nil
	case ParserStateDone:
		return 0, ErrorRequestAlreadyParsed
	default:
		return 0, ErrorUnknownParserState
	}
}

func (r *Request) done() bool {
	return r.ParserState == ParserStateDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		ParserState: ParserStateRequestLine,
		Headers:     headers.NewHeaders(),
	}

	buffer := make([]byte, bufferSize)
	readToIndex := 0

	for !request.done() {
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
				if !request.done() {
					return nil, ErrorIncompleteRequest
				}
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
		return nil, 0, ErrorRequestLineMalformed
	}

	if !slices.Contains(supportedHTTPMethods, parts[0]) {
		return nil, 0, ErrorHTTPMethodNotSupported
	}

	httpVersion := strings.TrimPrefix(parts[2], "HTTP/")
	if !slices.Contains(supportedHTTPVersions, httpVersion) {
		return nil, 0, ErrorHTTPVersionNotSupported
	}

	return &RequestLine{
		HTTPVersion:   httpVersion,
		RequestTarget: parts[1],
		Method:        parts[0],
	}, len(requestLine) + len(crlf), nil
}
