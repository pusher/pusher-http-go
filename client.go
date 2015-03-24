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

func (c *Client) trigger(channels []string, event string, _data interface{}, socketId string) (*BufferedEvents, error) {

	if len(channels) > 10 {
		return nil, errors.New("You cannot trigger on more than 10 channels at once")
	}

	if !channelsAreValid(channels) {
		return nil, errors.New("At least one of your channels' names are invalid")
	}

	payload, size_err := createTriggerPayload(channels, event, _data, socketId)

	if size_err != nil {
		return nil, size_err
	}

	path := "/apps/" + c.AppId + "/" + "events"
	u := createRequestUrl("POST", c.Host, path, c.Key, c.Secret, auth_timestamp(), payload, nil)
	response, responseErr := request("POST", u, payload)

	if responseErr != nil {
		return nil, responseErr
	}

	return unmarshalledBufferedEvents(response), nil
}

func (c *Client) Trigger(channels []string, event string, _data interface{}) (*BufferedEvents, error) {
	return c.trigger(channels, event, _data, "")
}

func (c *Client) TriggerExclusive(channels []string, event string, _data interface{}, socketId string) (*BufferedEvents, error) {
	return c.trigger(channels, event, _data, socketId)
}

func (c *Client) Channels(additionalQueries map[string]string) (*ChannelsList, error) {
	path := "/apps/" + c.AppId + "/channels"
	u := createRequestUrl("GET", c.Host, path, c.Key, c.Secret, auth_timestamp(), nil, additionalQueries)
	response, err := request("GET", u, nil)
	if err != nil {
		return nil, err
	}
	return unmarshalledChannelsList(response), nil
}

func (c *Client) Channel(name string, additionalQueries map[string]string) (*Channel, error) {
	path := "/apps/" + c.AppId + "/channels/" + name
	u := createRequestUrl("GET", c.Host, path, c.Key, c.Secret, auth_timestamp(), nil, additionalQueries)
	response, err := request("GET", u, nil)
	if err != nil {
		return nil, err
	}
	return unmarshalledChannel(response, name), nil
}

func (c *Client) GetChannelUsers(name string) (*Users, error) {
	path := "/apps/" + c.AppId + "/channels/" + name + "/users"
	u := createRequestUrl("GET", c.Host, path, c.Key, c.Secret, auth_timestamp(), nil, nil)
	response, err := request("GET", u, nil)
	if err != nil {
		return nil, err
	}
	return unmarshalledChannelUsers(response), nil
}

func (c *Client) AuthenticateChannel(_params []byte, member ...MemberData) string {

	channelName, socketId := parseAuthRequestParams(_params)
	stringToSign := strings.Join([]string{socketId, channelName}, ":")
	isPresenceChannel := strings.HasPrefix(channelName, "presence-")

	if isPresenceChannel {
		var presenceData MemberData
		if member != nil {
			presenceData = member[0]
		}
		return c.authenticatePresenceChannel(_params, stringToSign, presenceData)
	} else {
		return c.authenticatePrivateChannel(_params, stringToSign)
	}

}

func (c *Client) authenticatePrivateChannel(_params []byte, stringToSign string) string {
	_response := createAuthMap(c.Key, c.Secret, stringToSign)
	response, _ := json.Marshal(_response)
	return string(response)
}

func (c *Client) authenticatePresenceChannel(_params []byte, stringToSign string, presenceData MemberData) string {

	_jsonUserData, _ := json.Marshal(presenceData)
	jsonUserData := string(_jsonUserData)

	stringToSign = strings.Join([]string{stringToSign, jsonUserData}, ":")

	_response := createAuthMap(c.Key, c.Secret, stringToSign)
	_response["channel_data"] = jsonUserData
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
