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

	result := client.AuthenticateChannel(post_params)

	assert.Equal(t, expected, result)

}

func TestPresenceChannelAuthentication(t *testing.T) {
	client := Client{
		Key:    "278d425bdf160c739803",
		Secret: "7ad3773142a6692b25b8",
	}

	post_params := []byte("channel_name=presence-foobar&socket_id=1234.1234")

	presence_data := MemberData{UserId: 10, UserInfo: map[string]string{"name": "Mr. Pusher"}}

	expected := "{\"auth\":\"278d425bdf160c739803:afaed3695da2ffd16931f457e338e6c9f2921fa133ce7dac49f529792be6304c\",\"channel_data\":\"{\\\"user_id\\\":10,\\\"user_info\\\":{\\\"name\\\":\\\"Mr. Pusher\\\"}}\"}"

	result := client.AuthenticateChannel(post_params, presence_data)

	assert.Equal(t, expected, result)

}
