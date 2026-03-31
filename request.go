package pusher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

const (
	contentTypeHeaderKey   = "Content-Type"
	contentTypeHeaderValue = "application/json"
)

var headers = map[string]string{
	"Content-Type":     "application/json",
	"X-Pusher-Library": fmt.Sprintf("%s %s", libraryName, libraryVersion),
}

func request(ctx context.Context, client *http.Client, method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

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
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return responseBody, nil
	}
	return nil, fmt.Errorf("Status Code: %d - %s", response.StatusCode, responseBody)
}
