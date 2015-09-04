package pusher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseAuthRequestParamsNoSock(t *testing.T) {
	params := "channel_name=hello"
	_, _, result := parseAuthRequestParams([]byte(params))
	assert.Error(t, result)
	assert.EqualError(t, result, "Socket_id not found")
}

func TestParseAuthRequestParamsNoChan(t *testing.T) {
	params := "socket_id=45.3"
	_, _, result := parseAuthRequestParams([]byte(params))
	assert.Error(t, result)
	assert.EqualError(t, result, "Channel param not found")
}

func TestInvalidAuthParams(t *testing.T) {
	params := "%$@£$${}$£%|$^%$^|"
	_, _, result := parseAuthRequestParams([]byte(params))
	assert.Error(t, result)
}
