package headers

import (
	"slices"
	"strings"
	"unicode"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	splitReq := strings.Split(string(data), crlf)
	if len(splitReq) <= 1 {
		return 0, false, nil
	}

	fieldLine := splitReq[0]

	if fieldLine == "" {
		return len(crlf), true, nil
	}

	colonIndex := strings.Index(fieldLine, ":")
	if colonIndex == -1 {
		return 0, false, ErrorInvalidHeaderFormat
	}

	fieldName := strings.TrimLeftFunc(fieldLine[:colonIndex], func(r rune) bool {
		return unicode.IsSpace(r)
	})

	if unicode.IsSpace(rune(fieldName[len(fieldName)-1])) {
		return 0, false, ErrorInvalidFieldNameFormat
	}

	if !isValidToken([]byte(fieldName)) {
		return 0, false, ErrorInvalidFieldNameToken
	}

	fieldValue := strings.TrimSpace(fieldLine[colonIndex+1:])
	h.Set(fieldName, fieldValue)

	return len(fieldLine) + len(crlf), false, nil
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
