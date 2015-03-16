package pusher

import (
	"encoding/json"
	// "fmt"
	// "io"
	// "io/ioutil"
	"net/http"
	// "net/http/httputil"
	"net/url"
	"strings"
)

type Client struct {
	AppId, Key, Secret string
}

func (c *Client) trigger(channels []string, event string, _data interface{}, socket_id string) (error, string) {
	data, _ := json.Marshal(_data)

	payload, _ := json.Marshal(&Event{
		Name:     event,
		Channels: channels,
		Data:     string(data),
		SocketId: socket_id})

	path := "/apps/" + c.AppId + "/" + "events"

	u := Url{"POST", path, c.Key, c.Secret, payload, nil}

	err, response := Request("POST", u.generate(), payload)

	return err, string(response)
}

func (c *Client) Trigger(channels []string, event string, _data interface{}) (error, string) {
	return c.trigger(channels, event, _data, "")
}

func (c *Client) TriggerExclusive(channels []string, event string, _data interface{}, socket_id string) (error, string) {
	return c.trigger(channels, event, _data, socket_id)
}

func (c *Client) Channels(additional_queries map[string]string) (error, *ChannelsList) {
	path := "/apps/" + c.AppId + "/channels"
	u := Url{"GET", path, c.Key, c.Secret, nil, additional_queries}
	err, response := Request("GET", u.generate(), nil)

	channels := &ChannelsList{}
	json.Unmarshal(response, &channels)
	return err, channels
}

func (c *Client) Channel(name string, additional_queries map[string]string) (error, *Channel) {

	path := "/apps/" + c.AppId + "/channels/" + name

	u := Url{"GET", path, c.Key, c.Secret, nil, additional_queries}

	err, raw_channel_data := Request("GET", u.generate(), nil)

	channel := &Channel{Name: name}
	json.Unmarshal(raw_channel_data, &channel)
	return err, channel

}

func (c *Client) GetChannelUsers(name string) (error, *Users) {
	path := "/apps/" + c.AppId + "/channels/" + name + "/users"
	u := Url{"GET", path, c.Key, c.Secret, nil, nil}
	err, raw_users := Request("GET", u.generate(), nil)
	users := &Users{}
	json.Unmarshal(raw_users, &users)
	return err, users
}

func (c *Client) AuthenticateChannel(_params []byte, presence_data MemberData) string {
	params, _ := url.ParseQuery(string(_params))
	channel_name := params["channel_name"][0]
	socket_id := params["socket_id"][0]

	string_to_sign := socket_id + ":" + channel_name

	is_presence_channel := strings.HasPrefix(channel_name, "presence-")

	var json_user_data string
	_response := make(map[string]string)

	if is_presence_channel {
		_json_user_data, _ := json.Marshal(presence_data)
		json_user_data = string(_json_user_data)
		string_to_sign += ":" + json_user_data

		_response["channel_data"] = json_user_data
	}

	auth_signature := HMACSignature(string_to_sign, c.Secret)
	_response["auth"] = c.Key + ":" + auth_signature
	response, _ := json.Marshal(_response)

	return string(response)
}

func (c *Client) Webhook(header http.Header, body []byte) *Webhook {
	webhook := &Webhook{Key: c.Key, Secret: c.Secret, Header: header, RawBody: string(body)}
	json.Unmarshal(body, &webhook)
	return webhook
}
