package headers

import "errors"

var (
	ErrorInvalidHeaderFormat    = errors.New("invalid header format")
	ErrorInvalidFieldNameFormat = errors.New("invalid field name format")

	ErrorInvalidFieldNameToken = errors.New("invalid field name token")
)
