package pusher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrivateChannelAuthentication(t *testing.T) {
	client := Client{
		Key:    "278d425bdf160c739803",
		Secret: "7ad3773142a6692b25b8",
	}
	post_params := []byte("channel_name=private-foobar&socket_id=1234.1234")
	expected := "{\"auth\":\"278d425bdf160c739803:58df8b0c36d6982b82c3ecf6b4662e34fe8c25bba48f5369f135bf843651c3a4\"}"
	result, err := client.AuthenticatePrivateChannel(post_params)
	assert.Equal(t, expected, string(result))
	assert.NoError(t, err)
}

func TestPresenceChannelAuthentication(t *testing.T) {
	client := Client{
		Key:    "278d425bdf160c739803",
		Secret: "7ad3773142a6692b25b8",
	}
	post_params := []byte("channel_name=presence-foobar&socket_id=1234.1234")
	presence_data := MemberData{UserId: "10", UserInfo: map[string]string{"name": "Mr. Pusher"}}
	expected := "{\"auth\":\"278d425bdf160c739803:48dac51d2d7569e1e9c0f48c227d4b26f238fa68e5c0bb04222c966909c4f7c4\",\"channel_data\":\"{\\\"user_id\\\":\\\"10\\\",\\\"user_info\\\":{\\\"name\\\":\\\"Mr. Pusher\\\"}}\"}"
	result, err := client.AuthenticatePresenceChannel(post_params, presence_data)
	assert.Equal(t, expected, string(result))
	assert.NoError(t, err)
}
