package pusher

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var PusherURLRegex = regexp.MustCompile("^(http|https)://(.*):(.*)@(.*)/apps/([0-9]+)$")

type Client struct {
	AppId, Key, Secret, Host string
	Timeout                  time.Duration
}

func ClientFromURL(url string) (*Client, error) {
	matches := PusherURLRegex.FindStringSubmatch(url)
	if len(matches) == 0 {
		return nil, errors.New("No match found")
	}
	return &Client{Key: matches[2], Secret: matches[3], Host: matches[4], AppId: matches[5]}, nil
}

func ClientFromEnv(key string) (*Client, error) {
	url := os.Getenv(key)
	return ClientFromURL(url)
}

// triggerMulti
// socketId pointer to string

func (c *Client) Trigger(channel string, event string, data interface{}) (*BufferedEvents, error) {
	return c.trigger([]string{channel}, event, data, nil)
}

func (c *Client) TriggerMulti(channels []string, event string, data interface{}) (*BufferedEvents, error) {
	return c.trigger(channels, event, data, nil)
}

func (c *Client) TriggerExclusive(channel string, event string, data interface{}, socketID string) (*BufferedEvents, error) {
	return c.trigger([]string{channel}, event, data, &socketID)
}

func (c *Client) TriggerMultiExclusive(channels []string, event string, data interface{}, socketID string) (*BufferedEvents, error) {
	return c.trigger(channels, event, data, &socketID)
}

func (c *Client) trigger(channels []string, event string, data interface{}, socketId *string) (*BufferedEvents, error) {
	if len(channels) > 10 {
		return nil, errors.New("You cannot trigger on more than 10 channels at once")
	}

	if !channelsAreValid(channels) {
		return nil, errors.New("At least one of your channels' names are invalid")
	}

	payload, err := createTriggerPayload(channels, event, data, socketId)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/apps/%s/events", c.AppId)
	u := createRequestUrl("POST", c.Host, path, c.Key, c.Secret, authTimestamp(), payload, nil)
	response, err := request("POST", u, payload, c.Timeout)
	if err != nil {
		return nil, err
	}

	return unmarshalledBufferedEvents(response)
}

func (c *Client) Channels(additionalQueries map[string]string) (*ChannelsList, error) {
	path := fmt.Sprintf("/apps/%s/channels", c.AppId)
	u := createRequestUrl("GET", c.Host, path, c.Key, c.Secret, authTimestamp(), nil, additionalQueries)
	response, err := request("GET", u, nil, c.Timeout)
	if err != nil {
		return nil, err
	}
	return unmarshalledChannelsList(response)
}

func (c *Client) Channel(name string, additionalQueries map[string]string) (*Channel, error) {
	path := fmt.Sprintf("/apps/%s/channels/%s", c.AppId, name)
	u := createRequestUrl("GET", c.Host, path, c.Key, c.Secret, authTimestamp(), nil, additionalQueries)
	response, err := request("GET", u, nil, c.Timeout)
	if err != nil {
		return nil, err
	}
	return unmarshalledChannel(response, name)
}

func (c *Client) GetChannelUsers(name string) (*Users, error) {
	path := fmt.Sprintf("/apps/%s/channels/%s/users", c.AppId, name)
	u := createRequestUrl("GET", c.Host, path, c.Key, c.Secret, authTimestamp(), nil, nil)
	response, err := request("GET", u, nil, c.Timeout)
	if err != nil {
		return nil, err
	}
	return unmarshalledChannelUsers(response)
}

func (c *Client) AuthenticatePrivateChannel(params []byte) (response []byte, err error) {
	return c.authenticateChannel(params, nil)
}

func (c *Client) AuthenticatePresenceChannel(params []byte, member MemberData) (response []byte, err error) {
	return c.authenticateChannel(params, &member)
}

func (c *Client) authenticateChannel(params []byte, member *MemberData) (response []byte, err error) {
	channelName, socketId := parseAuthRequestParams(params)
	stringToSign := strings.Join([]string{socketId, channelName}, ":")

	var jsonUserData string

	if member != nil {
		var _jsonUserData []byte
		_jsonUserData, err = json.Marshal(member)
		if err != nil {
			return
		}

		jsonUserData = string(_jsonUserData)
		stringToSign = strings.Join([]string{stringToSign, jsonUserData}, ":")
	}

	_response := createAuthMap(c.Key, c.Secret, stringToSign)

	if member != nil {
		_response["channel_data"] = jsonUserData
	}

	response, err = json.Marshal(_response)
	return
}

func (c *Client) Webhook(header http.Header, body []byte) (*Webhook, error) {

	for _, token := range header["X-Pusher-Key"] {
		if token == c.Key && checkSignature(header.Get("X-Pusher-Signature"), string(body), c.Secret) {
			return unmarshalledWebhook(body)
		}
	}
	return nil, errors.New("Invalid webhook")
}
