package authentications

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrivateChannelRequestOffersSocketIDAndChannelNameToBeSigned(t *testing.T) {
	channel := &PrivateChannel{
		[]byte("channel_name=private-foobar&socket_id=1234.1234"),
	}
	unsigned, err := channel.StringToSign()
	assert.NoError(t, err)
	assert.Equal(t, "1234.1234:private-foobar", unsigned)
}

func TestPrivateChannelErrorsOnUnparseableQuery(t *testing.T) {
	channel := &PrivateChannel{
		[]byte(`<><><><><><><><>@<£>@3,@>%<£>%<^>%<^>^&<^>*<&>*&>&<*>^&<&>^<%>^<$>%4∞¢§∞§¶•¶`),
	}
	unsigned, err := channel.StringToSign()
	assert.Equal(t, "", unsigned)
	assert.EqualError(t, err, `invalid URL escape "%<\xc2"`)
}

func TestPrivateChannelRequestErrorsOnNoChannelParam(t *testing.T) {
	channel := &PrivateChannel{
		[]byte("socket_id=1234.1234"),
	}
	unsigned, err := channel.StringToSign()
	assert.Equal(t, "", unsigned)
	assert.EqualError(t, err, "[pusher-http-go]: Channel param not found")
}

func TestPrivateChannelRequestErrorsOnBlankChannelParam(t *testing.T) {
	channel := &PrivateChannel{
		[]byte("socket_id=1234.1234&channel_name="),
	}
	unsigned, err := channel.StringToSign()
	assert.Equal(t, "", unsigned)
	assert.EqualError(t, err, "[pusher-http-go]: Channel name cannot be blank")
}

func TestPrivateChannelRequestErrorsOnNoSocketID(t *testing.T) {
	channel := &PrivateChannel{
		[]byte("channel_name=private-yolo"),
	}
	unsigned, err := channel.StringToSign()
	assert.Equal(t, "", unsigned)
	assert.EqualError(t, err, "[pusher-http-go]: Socket_id not found")
}

func TestErrorOnInvalidSocketID(t *testing.T) {
	channel := &PrivateChannel{
		[]byte("channel_name=yolo&socket_id=fdgfd"),
	}
	_, err := channel.StringToSign()
	assert.EqualError(t, err, "[pusher-http-go]: socket_id invalid")
}

func TestPrivateChannelReturnsNilUserData(t *testing.T) {
	channel := &PrivateChannel{[]byte("yolo")}
	user, err := channel.UserData()
	assert.NoError(t, err)
	assert.Equal(t, "", user)
}
