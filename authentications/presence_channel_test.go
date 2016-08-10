package authentications

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockMember struct{}

func (m mockMember) UserData() (string, error) {
	return "{\"user_id\":\"10\",\"user_info\":{\"name\":\"Mr. Pusher\"}}", nil
}

type mockMemberReturnsError struct{}

func (m mockMemberReturnsError) UserData() (string, error) {
	return "", errors.New("Some JSON error")
}

func TestPresenceChannelOffersUserDataPlusSocketIDAndChannelNameToSign(t *testing.T) {
	presenceChannel := &PresenceChannel{
		Body:   []byte("channel_name=presence-foobar&socket_id=1234.1234"),
		Member: mockMember{},
	}
	unsigned, err := presenceChannel.StringToSign()
	expected :=
		"1234.1234:presence-foobar:{\"user_id\":\"10\",\"user_info\":{\"name\":\"Mr. Pusher\"}}"
	assert.Equal(t, expected, unsigned)
	assert.NoError(t, err)
}

func TestIfUserDataErrStringToSignErrors(t *testing.T) {
	presenceChannel := &PresenceChannel{
		Body:   []byte("channel_name=presence-foobar&socket_id=1234.1234"),
		Member: mockMemberReturnsError{},
	}
	_, err := presenceChannel.StringToSign()
	assert.EqualError(t, err, "Some JSON error")
}
