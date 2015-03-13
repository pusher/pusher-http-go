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

	u := Url{"POST", path, c.Key, c.Secret, payload, nil}

	err, response := Request("POST", u.generate(), payload)

	return err, string(response)
}

func (c *Client) Channels(additional_queries map[string]string) (error, string) {

	path := "/apps/" + c.AppId + "/channels"
	u := Url{"GET", path, c.Key, c.Secret, nil, additional_queries}
	err, response := Request("GET", u.generate(), nil)
	return err, string(response)
}

func (c *Client) Channel(name string, additional_queries map[string]string) (error, *Channel) {

	path := "/apps/" + c.AppId + "/channels/" + name

	u := Url{"GET", path, c.Key, c.Secret, nil, additional_queries}

	err, raw_channel_data := Request("GET", u.generate(), nil)

	channel := &Channel{Name: name, Client: *c}
	json.Unmarshal(raw_channel_data, &channel)
	return err, channel

}
