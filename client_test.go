package pusher

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestTrigger(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		res.WriteHeader(200)
		fmt.Fprintf(res, "{}")

	}))

	defer server.Close()

	u, _ := url.Parse(server.URL)

	client := Client{"id", "key", "secret", u.Host}

	err, _ := client.Trigger([]string{"test_channel"}, "test", "yolo")

	assert.NoError(t, err)

}
