package requests

import (
	"bytes"
	"fmt"
	"github.com/pusher/pusher/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Request struct {
	Method      string
	PathPattern string
}

func (r *Request) Do(client *http.Client, u *url.URL, payload []byte) (responseBody []byte, err error) {
	var (
		httpResponse *http.Response
		httpRequest  *http.Request
	)

	defer func() {
		if httpResponse != nil {
			httpResponse.Body.Close()
		}
	}()

	if httpRequest, err = http.NewRequest(r.Method, u.String(), bytes.NewReader(payload)); err != nil {
		return
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	if httpResponse, err = client.Do(httpRequest); err != nil {
		return
	}

	if responseBody, err = ioutil.ReadAll(httpResponse.Body); err != nil {
		return
	}

	if httpResponse.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Status Code: %s - %s", httpResponse.Status, string(responseBody)))
	}
	return
}
