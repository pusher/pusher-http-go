package validate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSocketIDValidation(t *testing.T) {
	invalidSocketID := "12341234"
	err := SocketID(&invalidSocketID)
	assert.EqualError(t, err, "[pusher-http-go]: socket_id invalid")
}

func TestNoSocketIDNoError(t *testing.T) {
	err := SocketID(nil)
	assert.NoError(t, err)
}
