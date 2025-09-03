// Package headers provides functionality to parse and manage HTTP headers.
package headers

import (
	"bytes"
	"slices"
	"strings"
	"unicode"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIdx := bytes.Index(data, []byte(crlf))
	if crlfIdx == -1 {
		// We don't have a full line yet
		return 0, false, nil
	}

	// Check if the line is just CRLF indicating the end of headers
	if crlfIdx == 0 {
		return len(crlf), true, nil
	}

	// We have at least one full line to process
	splitReq := bytes.SplitN(data[:crlfIdx], []byte(":"), 2)
	if len(splitReq) != 2 {
		return 0, false, nil
	}

	fieldName := splitReq[0]
	if unicode.IsSpace(rune(fieldName[len(fieldName)-1])) {
		return 0, false, ErrorInvalidFieldNameFormat
	}

	fieldName = bytes.TrimSpace(fieldName)
	if !isValidToken([]byte(fieldName)) {
		return 0, false, ErrorInvalidFieldNameToken
	}

	fieldValue := bytes.TrimSpace(splitReq[1])

	h.Set(string(fieldName), string(fieldValue))

	return crlfIdx + len(crlf), false, nil
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func isValidToken(data []byte) bool {
	for _, token := range data {
		if (token < 'a' || token > 'z') &&
			(token < 'A' || token > 'Z') &&
			(token < '0' || token > '9') &&
			!slices.Contains(tokenChars, token) {
			return false
		}
	}

	return true
}

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Get(key string) (string, bool) {
	value, ok := h[strings.ToLower(key)]
	return value, ok
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)

	existingValue, ok := h[key]
	if ok {
		h[key] = existingValue + ", " + value
	} else {
		h[key] = value
	}
}
