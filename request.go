package pusher

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func request(method, url string, body []byte) ([]byte, error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	return process_response(resp.StatusCode, resp_body)
}

func process_response(status int, resp_body []byte) ([]byte, error) {
	if status == 200 {
		return resp_body, nil
	} else {
		message := fmt.Sprintf("Status Code: %s - %s", strconv.Itoa(status), string(resp_body))
		err := errors.New(message)
		return nil, err
	}
}
