package requests

import (
	"bytes"
	"fmt"
	"github.com/pusher/pusher-http-go/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Request struct {
	Method      string
	PathPattern string
}

func (req *Request) Do(client *http.Client, u *url.URL, payload []byte) ([]byte, error) {
	httpRequest, err := http.NewRequest(req.Method, u.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	responseBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusAccepted {
		return nil, errors.New(fmt.Sprintf("Status Code: %s - %s", httpResponse.Status, string(responseBody)))
	}

	return responseBody, nil
}
