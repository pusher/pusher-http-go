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

	path := "/apps/" + c.AppId + "/" + "events"

	q := Query{"POST", path, c.Key, c.Secret, payload, nil}

	return c.post(q.generate(), payload)
}

func (c *Client) Channels(additional_queries map[string]string) (error, string) {

	path := "/apps/" + c.AppId + "/channels"
	q := Query{"GET", path, c.Key, c.Secret, nil, additional_queries}
	return c.get(q.generate(), nil)
}

func (c *Client) Channel(name string, additional_queries map[string]string) (error, string) {

	path := "/apps/" + c.AppId + "/channels/" + name

	q := Query{"GET", path, c.Key, c.Secret, nil, additional_queries}

	return c.get(q.generate(), nil)
}

func (c *Client) post(url string, body []byte) (error, string) {
	request := &Request{"POST", url, body}
	return request.send()
}

func (c *Client) get(url string, body []byte) (error, string) {
	request := &Request{"GET", url, body}
	return request.send()
}
