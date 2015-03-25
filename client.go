package pusher

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	u "net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var pusherPathRegex = regexp.MustCompile("^/apps/([0-9]+)$")

// Client to the HTTP API of Pusher
type Client struct {
	AppId   string
	Key     string
	Secret  string
	Host    string        // host or host:port pair
	Secure  bool          // true for HTTPS
	Timeout time.Duration // Request timeout for HTTP requests
}

// Client constructor from a specially crafted URL
//
// Eg: ClientFromURL("http://feaf18a411d3cb9216ee:fec81108d90e1898e17a@api.pusherapp.com/apps/104060")
func ClientFromURL(url string) (*Client, error) {
	url2, err := u.Parse(url)
	if err != nil {
		return nil, err
	}

	c := Client{
		Host: url2.Host,
	}

	matches := pusherPathRegex.FindStringSubmatch(url2.Path)
	if len(matches) == 0 {
		return nil, errors.New("No app ID found")
	}
	c.AppId = matches[1]

	if url2.User == nil {
		return nil, errors.New("Missing <key>:<secret>")
	}
	c.Key = url2.User.Username()
	var isSet bool
	c.Secret, isSet = url2.User.Password()
	if !isSet {
		return nil, errors.New("Missing <secret>")
	}

	if url2.Scheme == "https" {
		c.Secure = true
	}

	return &c, nil
}

// Client constructor for an environment variable (like Heroku).
//
// Eg: ClientFromEnv("PUSHER_URL")
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
	u := createRequestUrl("POST", c.Host, path, c.Key, c.Secret, authTimestamp(), c.Secure, payload, nil)
	response, err := request("POST", u, payload, c.Timeout)
	if err != nil {
		return nil, err
	}

	return unmarshalledBufferedEvents(response)
}

func (c *Client) Channels(additionalQueries map[string]string) (*ChannelsList, error) {
	path := fmt.Sprintf("/apps/%s/channels", c.AppId)
	u := createRequestUrl("GET", c.Host, path, c.Key, c.Secret, authTimestamp(), c.Secure, nil, additionalQueries)
	response, err := request("GET", u, nil, c.Timeout)
	if err != nil {
		return nil, err
	}
	return unmarshalledChannelsList(response)
}

func (c *Client) Channel(name string, additionalQueries map[string]string) (*Channel, error) {
	path := fmt.Sprintf("/apps/%s/channels/%s", c.AppId, name)
	u := createRequestUrl("GET", c.Host, path, c.Key, c.Secret, authTimestamp(), c.Secure, nil, additionalQueries)
	response, err := request("GET", u, nil, c.Timeout)
	if err != nil {
		return nil, err
	}
	return unmarshalledChannel(response, name)
}

func (c *Client) GetChannelUsers(name string) (*Users, error) {
	path := fmt.Sprintf("/apps/%s/channels/%s/users", c.AppId, name)
	u := createRequestUrl("GET", c.Host, path, c.Key, c.Secret, authTimestamp(), c.Secure, nil, nil)
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
	channelName, socketId, err := parseAuthRequestParams(params)

	if err != nil {
		return
	}

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
