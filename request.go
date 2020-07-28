package pusher

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	contentTypeHeaderKey   = "Content-Type"
	contentTypeHeaderValue = "application/json"
)

var headers = map[string]string{
	"Content-Type":     "application/json",
	"X-Pusher-Library": fmt.Sprintf("%s %s", libraryName, libraryVersion),
}

// change timeout to time.Duration
func request(client *http.Client, method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	for key, val := range headers {
		req.Header.Set(http.CanonicalHeaderKey(key), val)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return processResponse(resp)
}

func processResponse(response *http.Response) ([]byte, error) {
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return responseBody, nil
	}
	message := fmt.Sprintf("Status Code: %s - %s", strconv.Itoa(response.StatusCode), string(responseBody))
	err = errors.New(message)
	return nil, err
}
