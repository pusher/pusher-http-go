package pusher

import (
	"encoding/json"
	// "fmt"
)

type Client struct {
	AppId, Key, Secret string
}

func (c *Client) Trigger(channels []string, event string, _data map[string]string) (error, string) {
	data, _ := json.Marshal(_data)

	payload, _ := json.Marshal(&EventBody{
		Name:     event,
		Channels: channels,
		Data:     string(data)})

	q := Query{"POST", c.path("events"), c.Key, c.Secret, payload}

	return c.post(q.generate(), payload)
}

func (c *Client) Channels() (error, string) {
	q := Query{"GET", c.path("channels"), c.Key, c.Secret, nil}
	return c.get(q.generate(), nil)
}

func (c *Client) path(resource string) string {
	return "/apps/" + c.AppId + "/" + resource
}

func (c *Client) post(url string, body []byte) (error, string) {
	request := &Request{"POST", url, body}
	return request.send()
}

func (c *Client) get(url string, body []byte) (error, string) {
	request := &Request{"GET", url, body}
	return request.send()
}
