package pusher

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
	"strconv"
)

// change timeout to time.Duration
func request(client *http.Client, method, url string, body []byte, logger *zerolog.Logger) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if logger != nil {
			logger.Error().Err(err).Msgf("cannot do http request %+v to %s", req, url)
		}
		return nil, err
	}
	defer resp.Body.Close()
	return processResponse(resp, logger)
}

func processResponse(response *http.Response, logger *zerolog.Logger) ([]byte, error) {
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		if logger != nil {
			logger.Error().Err(err).Msgf("cannot read response body from %+v", response)
		}
		return nil, err
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return responseBody, nil
	}
	message := fmt.Sprintf("Status Code: %s - %s", strconv.Itoa(response.StatusCode), string(responseBody))
	err = errors.New(message)
	if logger != nil {
		logger.Error().Err(err).Msg(message)
	}
	return nil, err
}
