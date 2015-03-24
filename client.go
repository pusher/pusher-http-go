package pusher

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Client struct {
	AppId, Key, Secret, Host string
}

func (c *Client) trigger(channels []string, event string, _data interface{}, socket_id string) (error, *BufferedEvents) {

	if len(channels) > 10 {
		return errors.New("You cannot trigger on more than 10 channels at once"), nil
	}

	if !channelsAreValid(channels) {
		return errors.New("At least one of your channels' names are invalid"), nil
	}

	payload, size_err := createTriggerPayload(channels, event, _data, socket_id)

	if size_err != nil {
		return size_err, nil
	}

	path := "/apps/" + c.AppId + "/" + "events"
	u := createRequestUrl("POST", c.Host, path, c.Key, c.Secret, auth_timestamp(), payload, nil)
	response_err, response := request("POST", u, payload)

	if response_err != nil {
		return response_err, nil
	}

	return nil, unmarshalledBufferedEvents(response)
}

func (c *Client) Trigger(channels []string, event string, _data interface{}) (error, *BufferedEvents) {
	return c.trigger(channels, event, _data, "")
}

func (c *Client) TriggerExclusive(channels []string, event string, _data interface{}, socket_id string) (error, *BufferedEvents) {
	return c.trigger(channels, event, _data, socket_id)
}

func (c *Client) Channels(additional_queries map[string]string) (error, *ChannelsList) {
	path := "/apps/" + c.AppId + "/channels"
	u := createRequestUrl("GET", c.Host, path, c.Key, c.Secret, auth_timestamp(), nil, additional_queries)
	err, response := request("GET", u, nil)
	if err != nil {
		return err, nil
	}
	return nil, unmarshalledChannelsList(response)
}

func (c *Client) Channel(name string, additional_queries map[string]string) (error, *Channel) {
	path := "/apps/" + c.AppId + "/channels/" + name
	u := createRequestUrl("GET", c.Host, path, c.Key, c.Secret, auth_timestamp(), nil, additional_queries)
	err, response := request("GET", u, nil)
	if err != nil {
		return err, nil
	}
	return nil, unmarshalledChannel(response, name)
}

func (c *Client) GetChannelUsers(name string) (error, *Users) {
	path := "/apps/" + c.AppId + "/channels/" + name + "/users"
	u := createRequestUrl("GET", c.Host, path, c.Key, c.Secret, auth_timestamp(), nil, nil)
	err, response := request("GET", u, nil)
	if err != nil {
		return err, nil
	}
	return nil, unmarshalledChannelUsers(response)
}

func (c *Client) AuthenticateChannel(_params []byte, member ...MemberData) string {

	channel_name, socket_id := parseAuthRequestParams(_params)
	string_to_sign := strings.Join([]string{socket_id, channel_name}, ":")
	is_presence_channel := strings.HasPrefix(channel_name, "presence-")

	if is_presence_channel {
		var presence_data MemberData
		if member != nil {
			presence_data = member[0]
		}
		return c.authenticatePresenceChannel(_params, string_to_sign, presence_data)
	} else {
		return c.authenticatePrivateChannel(_params, string_to_sign)
	}

}

func (c *Client) authenticatePrivateChannel(_params []byte, string_to_sign string) string {
	_response := createAuthMap(c.Key, c.Secret, string_to_sign)
	response, _ := json.Marshal(_response)
	return string(response)
}

func (c *Client) authenticatePresenceChannel(_params []byte, string_to_sign string, presence_data MemberData) string {

	_json_user_data, _ := json.Marshal(presence_data)
	json_user_data := string(_json_user_data)

	string_to_sign = strings.Join([]string{string_to_sign, json_user_data}, ":")

	_response := createAuthMap(c.Key, c.Secret, string_to_sign)
	_response["channel_data"] = json_user_data
	response, _ := json.Marshal(_response)
	return string(response)
}

func (c *Client) Webhook(header http.Header, body []byte) (*Webhook, error) {
	for _, token := range header["X-Pusher-Key"] {
		if token == c.Key && checkSignature(header["X-Pusher-Signature"][0], string(body), c.Secret) {
			return unmarshalledWebhook(body), nil
		}
	}
	return nil, errors.New("Invalid webhook")
}
