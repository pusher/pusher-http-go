package requests

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRequestSentWithContentTypeJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, http.MethodOptions, req.Method)
		var (
			body []byte
			err  error
		)
		body, err = ioutil.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, []byte("payload"), body)
		assert.Equal(t, "application/json", req.Header["Content-Type"][0])
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(":)"))
	}))
	defer server.Close()

	req := &Request{
		Method:      http.MethodOptions,
		PathPattern: "yolo",
	}

	u, _ := url.Parse(server.URL)
	res, err := req.Do(http.DefaultClient, u, []byte("payload"))
	assert.Equal(t, []byte(":)"), res)
	assert.NoError(t, err)
}

func TestHandlingOfNon200StatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("You're a disgrace."))
	}))
	defer server.Close()

	req := &Request{
		Method:      http.MethodDelete,
		PathPattern: "bla",
	}
	u, _ := url.Parse(server.URL)
	_, err := req.Do(http.DefaultClient, u, nil)
	assert.EqualError(t, err, "[pusher-http-go]: Status Code: 403 Forbidden - You're a disgrace.")
}
