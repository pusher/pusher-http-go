package pusher

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	StatusCode int
	Body       []byte
}

func request(method, url string, body []byte, timeout int) ([]byte, error) {

	responseChannel := make(chan Response)
	httpClientError := make(chan error)
	client := &http.Client{}

	if timeout == 0 {
		timeout = 5
	}

	go func() {
		req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			httpClientError <- err
		}

		responseBody, _ := ioutil.ReadAll(resp.Body)
		response := Response{resp.StatusCode, responseBody}
		responseChannel <- response
	}()

	select {
	case response := <-responseChannel:
		return processResponse(response)
	case err := <-httpClientError:
		return nil, err
	case <-time.After(time.Duration(timeout) * time.Second):
		return nil, errors.New("The server was taking too long")
	}

}

func processResponse(response Response) ([]byte, error) {
	if response.StatusCode == 200 {
		return response.Body, nil
	} else {
		message := fmt.Sprintf("Status Code: %s - %s", strconv.Itoa(response.StatusCode), string(response.Body))
		err := errors.New(message)
		return nil, err
	}
}
