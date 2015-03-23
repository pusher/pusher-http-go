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
		assert.NoError(t, err)

	}))
	defer server.Close()
	u, _ := url.Parse(server.URL)
	client := Client{"id", "key", "secret", u.Host}
	err, _ := client.Trigger([]string{"test_channel"}, "test", "yolo")
	assert.NoError(t, err)
}

func Test400ResponseHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(400)
		fmt.Fprintf(res, "Cannot retrieve the user count unless the channel is a presence channel")

	}))

	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := Client{"id", "key", "secret", u.Host}

	channelParams := map[string]string{"info": "user_count,subscription_count"}
	err, _ := client.Channel("this_is_not_a_presence_channel", channelParams)

	assert.Error(t, err)

	assert.EqualError(t, err, "Status Code: 400 - Cannot retrieve the user count unless the channel is a presence channel")

}
