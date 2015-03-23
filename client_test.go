package pusher

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestTriggerSuccessCase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		fmt.Fprintf(res, "{}")
		assert.Equal(t, "POST", req.Method)

		expected_body := "{\"name\":\"test\",\"channels\":[\"test_channel\"],\"data\":\"\\\"yolo\\\"\",\"socket_id\":\"\"}"
		actual_body, err := ioutil.ReadAll(req.Body)
		assert.Equal(t, expected_body, string(actual_body))
		assert.Equal(t, "application/json", req.Header["Content-Type"][0])
		assert.NoError(t, err)

	}))
	defer server.Close()
	u, _ := url.Parse(server.URL)
	client := Client{"id", "key", "secret", u.Host}
	err, _ := client.Trigger([]string{"test_channel"}, "test", "yolo")
	assert.NoError(t, err)
}

func TestErrorResponseHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(400)
		fmt.Fprintf(res, "Cannot retrieve the user count unless the channel is a presence channel")

	}))

	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := Client{"id", "key", "secret", u.Host}

	channelParams := map[string]string{"info": "user_count,subscription_count"}
	err, channel := client.Channel("this_is_not_a_presence_channel", channelParams)

	assert.Error(t, err)
	assert.EqualError(t, err, "Status Code: 400 - Cannot retrieve the user count unless the channel is a presence channel")
	assert.Nil(t, channel)
}

func TestChannelLengthValidation(t *testing.T) {
	channels := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}

	client := Client{AppId: "id", Key: "key", Secret: "secret"}
	err, res := client.Trigger(channels, "yolo", "woot")

	assert.EqualError(t, err, "You cannot trigger on more than 10 channels at once")
	assert.Equal(t, "", res)
}

func TestChannelFormatValidation(t *testing.T) {
	channel1 := "w000^$$£@@@"

	channel2 := "lllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll"

	client := Client{AppId: "id", Key: "key", Secret: "secret"}
	err1, res1 := client.Trigger([]string{channel1}, "yolo", "w00t")

	err2, res2 := client.Trigger([]string{channel2}, "yolo", "not 19 forever")

	assert.EqualError(t, err1, "At least one of your channels' names are invalid")
	assert.Equal(t, "", res1)

	assert.EqualError(t, err2, "At least one of your channels' names are invalid")
	assert.Equal(t, "", res2)

}
