package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestLineParse(t *testing.T) {
	assert.Equal(t, "jyu.dev", "jyu.dev")
}
