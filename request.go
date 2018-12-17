package pusher

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// change timeout to time.Duration
func request(client *http.Client, method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
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
