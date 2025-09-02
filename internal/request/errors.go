package request

import "errors"

var (
	ErrorRequestAlreadyParsed = errors.New("request already fully parsed")
	ErrorUnknownParserState   = errors.New("unknown/unhandled parser state")

	ErrorRequestLineMalformed = errors.New("request line malformed")

	ErrorHTTPMethodNotSupported  = errors.New("http method not supported")
	ErrorHTTPVersionNotSupported = errors.New("http version not supported")
)
