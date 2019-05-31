package pusher

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var pusherPathRegex = regexp.MustCompile("^/apps/([0-9]+)$")
var maxTriggerableChannels = 100

/*
Client to the HTTP API of Pusher.

There easiest way to configure the library is by creating a `Pusher` instance:

    client := pusher.Client{
      AppID: "your_app_id",
      Key: "your_app_key",
      Secret: "your_app_secret",
    }

To ensure requests occur over HTTPS, set the `Encrypted` property of a
`pusher.Client` to `true`.

    client.Secure = true // false by default

If you wish to set a time-limit for each HTTP request, set the `Timeout`
property to an instance of `time.Duration`, for example:

    client.Timeout = time.Second * 3 // 5 seconds by default

Changing the `pusher.Client`'s `Host` property will make sure requests are sent
to your specified host.

    client.Host = "foo.bar.com" // by default this is "api.pusherapp.com".

*/
type Client struct {
	AppID               string
	Key                 string
	Secret              string
	Host                string // host or host:port pair
	Secure              bool   // true for HTTPS
	Cluster             string
	HTTPClient          *http.Client
	EncryptionMasterKey string //for E2E
}

/*
ClientFromURL allows client instantiation from a specially-crafted Pusher URL.

	c := pusher.ClientFromURL("http://key:secret@api.pusherapp.com/apps/app_id")

*/
func ClientFromURL(serverURL string) (*Client, error) {
	url2, err := url.Parse(serverURL)
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
	c.AppID = matches[1]

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

/*
ClientFromEnv allows instantiation of a client from an environment variable.
This is particularly relevant if you are using Pusher as a Heroku add-on,
which stores credentials in a `"PUSHER_URL"` environment variable. For example:
	client := pusher.ClientFromEnv("PUSHER_URL")

*/
func ClientFromEnv(key string) (*Client, error) {
	url := os.Getenv(key)
	return ClientFromURL(url)
}

/*
Returns the underlying HTTP client.
Useful to set custom properties to it.
*/
func (c *Client) requestClient() *http.Client {
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{Timeout: time.Second * 5}
	}

	return c.HTTPClient
}

func (c *Client) request(method, url string, body []byte) ([]byte, error) {
	return request(c.requestClient(), method, url, body)
}

/*
Trigger triggers an event to the Pusher API.
It is possible to trigger an event on one or more channels. Channel names can
contain only characters which are alphanumeric, `_` or `-`` and have
to be at most 200 characters long. Event name can be at most 200 characters long too.

Pass in the channel's name, the event's name, and a data payload. The data payload must
be marshallable into JSON.

	data := map[string]string{"hello": "world"}
	client.Trigger("greeting_channel", "say_hello", data)

*/
func (c *Client) Trigger(channel string, eventName string, data interface{}) error {
	return c.trigger([]string{channel}, eventName, data, nil)
}

/*
TriggerMulti is the same as `client.Trigger`, except one passes in a slice of
`channels` as the first parameter. The maximum length of channels is 100.
	client.TriggerMulti([]string{"a_channel", "another_channel"}, "event", data)
*/
func (c *Client) TriggerMulti(channels []string, eventName string, data interface{}) error {
	return c.trigger(channels, eventName, data, nil)
}

/*
TriggerExclusive triggers an event excluding a recipient whose connection has
the `socket_id` you specify here from receiving the event.
You can read more here: http://pusher.com/docs/duplicates.
	client.TriggerExclusive("a_channel", "event", data, "123.12")
*/
func (c *Client) TriggerExclusive(channel string, eventName string, data interface{}, socketID string) error {
	return c.trigger([]string{channel}, eventName, data, &socketID)
}

/*
TriggerMultiExclusive triggers an event to multiple channels excluding a
recipient whose connection has the `socket_id` you specify here from receiving
the event on any of the channels.
	client.TriggerMultiExclusive([]string{"a_channel", "another_channel"}, "event", data, "123.12")
*/
func (c *Client) TriggerMultiExclusive(channels []string, eventName string, data interface{}, socketID string) error {
	return c.trigger(channels, eventName, data, &socketID)
}

func (c *Client) trigger(channels []string, eventName string, data interface{}, socketID *string) error {
	hasEncryptedChannel := false
	for _, channel := range channels {
		if isEncryptedChannel(channel) {
			hasEncryptedChannel = true
		}
	}
	if len(channels) > maxTriggerableChannels {
		return fmt.Errorf("You cannot trigger on more than %d channels at once", maxTriggerableChannels)
	}
	if hasEncryptedChannel && len(channels) > 1 {
		// For rationale, see limitations of end-to-end encryption in the README
		return errors.New("You cannot trigger to multiple channels when using encrypted channels")

	}
	if !channelsAreValid(channels) {
		return errors.New("At least one of your channels' names are invalid")
	}
	if hasEncryptedChannel && !validEncryptionKey(c.EncryptionMasterKey) {
		return errors.New("Your encryptionMasterKey is not of the correct format")
	}
	if err := validateSocketID(socketID); err != nil {
		return err
	}
	payload, err := encodeTriggerBody(channels, eventName, data, socketID, c.EncryptionMasterKey)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/apps/%s/events", c.AppID)
	triggerURL, err := createRequestURL("POST", c.Host, path, c.Key, c.Secret, authTimestamp(), c.Secure, payload, nil, c.Cluster)
	if err != nil {
		return err
	}
	_, err = c.request("POST", triggerURL, payload)

	return err
}

/*
Event stores all the data for one Event that can be triggered.
*/
type Event struct {
	Channel  string
	Name     string
	Data     interface{}
	SocketID *string
}

/*
TriggerBatch triggers multiple events on multiple channels in a single call:
    client.TriggerBatch([]pusher.Event{
	    { Channel: "donut-1", Name: "ev1", Data: "d1" },
	    { Channel: "private-encrypted-secretdonut", Name: "ev2", Data: "d2" },
    })
*/
func (c *Client) TriggerBatch(batch []Event) error {
	hasEncryptedChannel := false
	// validate every channel name and every sockedID (if present) in batch
	for _, event := range batch {
		if !validChannel(event.Channel) {
			return fmt.Errorf("The channel named %s has a non-valid name", event.Channel)
		}
		if err := validateSocketID(event.SocketID); err != nil {
			return err
		}
		if isEncryptedChannel(event.Channel) {
			hasEncryptedChannel = true
		}
	}
	if hasEncryptedChannel {
		// validate EncryptionMasterKey
		if !validEncryptionKey(c.EncryptionMasterKey) {
			return errors.New("Your encryptionMasterKey is not of the correct format")
		}
	}
	payload, err := encodeTriggerBatchBody(batch, c.EncryptionMasterKey)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/apps/%s/batch_events", c.AppID)
	triggerURL, err := createRequestURL("POST", c.Host, path, c.Key, c.Secret, authTimestamp(), c.Secure, payload, nil, c.Cluster)
	if err != nil {
		return err
	}
	_, err = c.request("POST", triggerURL, payload)
	return err
}

/*
Channels returns a list of all the channels in an application. The parameter
`additionalQueries` is a map with query options. A key with `"filter_by_prefix"`
will filter the returned channels. To get number of users subscribed to a
presence-channel, specify an `"info"` key with value `"user_count"`. Pass in
`nil` if you do not wish to specify any query attributes

    channelsParams := map[string]string{
        "filter_by_prefix": "presence-",
        "info":             "user_count",
    }

    channels, err := client.Channels(channelsParams)

    //channels=> &{Channels:map[presence-chatroom:{UserCount:4} presence-notifications:{UserCount:31}  ]}

*/
func (c *Client) Channels(additionalQueries map[string]string) (*ChannelsList, error) {
	path := fmt.Sprintf("/apps/%s/channels", c.AppID)
	u, err := createRequestURL("GET", c.Host, path, c.Key, c.Secret, authTimestamp(), c.Secure, nil, additionalQueries, c.Cluster)
	if err != nil {
		return nil, err
	}
	response, err := c.request("GET", u, nil)
	if err != nil {
		return nil, err
	}
	return unmarshalledChannelsList(response)
}

/*
Channel allows you to get the state of a single channel. The parameter
`additionalQueries` is a map with query options. An `"info"` key can have
comma-separated vales of `"user_count"`, for presence-channels, and
`"subscription_count"`, for all-channels. Note that the subscription count is
not allowed by default. Please contact us at http://support.pusher.com if you
wish to enable this. Pass in `nil` if you do not wish to specify any query
attributes.

    channelParams := map[string]string{
        "info": "user_count,subscription_count",
    }

    channel, err := client.Channel("presence-chatroom", channelParams)

    //channel=> &{Name:presence-chatroom Occupied:true UserCount:42 SubscriptionCount:42}
*/
func (c *Client) Channel(name string, additionalQueries map[string]string) (*Channel, error) {
	path := fmt.Sprintf("/apps/%s/channels/%s", c.AppID, name)
	u, err := createRequestURL("GET", c.Host, path, c.Key, c.Secret, authTimestamp(), c.Secure, nil, additionalQueries, c.Cluster)
	if err != nil {
		return nil, err
	}
	response, err := c.request("GET", u, nil)
	if err != nil {
		return nil, err
	}
	return unmarshalledChannel(response, name)
}

/*
GetChannelUsers returns a list of users in a presence-channel by passing to this
method the channel name.

    users, err := client.GetChannelUsers("presence-chatroom")

    //users=> &{List:[{ID:13} {ID:90}]}

*/
func (c *Client) GetChannelUsers(name string) (*Users, error) {
	path := fmt.Sprintf("/apps/%s/channels/%s/users", c.AppID, name)
	u, err := createRequestURL("GET", c.Host, path, c.Key, c.Secret, authTimestamp(), c.Secure, nil, nil, c.Cluster)
	if err != nil {
		return nil, err
	}
	response, err := c.request("GET", u, nil)
	if err != nil {
		return nil, err
	}
	return unmarshalledChannelUsers(response)
}

/*
AuthenticatePrivateChannel allows you to authenticate a users subscription to a
private channel. It returns authentication signature to send back to the client
and authorize them.

For more information see our docs: http://pusher.com/docs/authenticating_users.

This is an example of authenticating a private-channel, using the built-in
Golang HTTP library to start a server.

In order to authorize a client, one must read the response into type `[]byte`
and pass it in. This will return a signature in the form of a `[]byte` for you
to send back to the client.

    func pusherAuth(res http.ResponseWriter, req *http.Request) {

        params, _ := ioutil.ReadAll(req.Body)
        response, err := client.AuthenticatePrivateChannel(params)

        if err != nil {
            panic(err)
        }

        fmt.Fprintf(res, string(response))

    }

    func main() {
        http.HandleFunc("/pusher/auth", pusherAuth)
        http.ListenAndServe(":5000", nil)
    }

*/
func (c *Client) AuthenticatePrivateChannel(params []byte) (response []byte, err error) {
	return c.authenticateChannel(params, nil)
}

/*
AuthenticatePresenceChannel allows you to authenticate a users subscription to a
presence channel. It returns authentication signature to send back to the client
and authorize them. In order to identify a user, clients are sent a user_id and,
optionally, custom data.

In this library, one does this by passing a `pusher.MemberData` instance.

	params, _ := ioutil.ReadAll(req.Body)

	presenceData := pusher.MemberData{
		UserID: "1",
		UserInfo: map[string]string{
			"twitter": "jamiepatel",
		},
	}

	response, err := client.AuthenticatePresenceChannel(params, presenceData)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(res, response)

*/
func (c *Client) AuthenticatePresenceChannel(params []byte, member MemberData) (response []byte, err error) {
	return c.authenticateChannel(params, &member)
}

func (c *Client) authenticateChannel(params []byte, member *MemberData) (response []byte, err error) {

	channelName, socketID, err := parseAuthRequestParams(params)
	if err != nil {
		return
	}

	if err = validateSocketID(&socketID); err != nil {
		return
	}

	stringToSign := strings.Join([]string{socketID, channelName}, ":")

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

	var _response map[string]string

	if isEncryptedChannel(channelName) {
		sharedSecret := generateSharedSecret(channelName, c.EncryptionMasterKey)
		sharedSecretB64 := base64.StdEncoding.EncodeToString(sharedSecret[:])
		_response = createAuthMap(c.Key, c.Secret, stringToSign, sharedSecretB64)
	} else {
		_response = createAuthMap(c.Key, c.Secret, stringToSign, "")
	}

	if member != nil {
		_response["channel_data"] = jsonUserData
	}

	response, err = json.Marshal(_response)
	return
}

/*
Webhook allows you to check that a Webhook you receive is indeed from Pusher, by
checking the token and authentication signature in the header of the request. On
your dashboard at http://app.pusher.com, you can set up webhooks to POST a
payload to your server after certain events. Such events include channels being
occupied or vacated, members being added or removed in presence-channels, or
after client-originated events. For more information see
https://pusher.com/docs/webhooks.



If the webhook is valid, a `*pusher.Webhook* will be returned, and the `err`
value will be nil. If it is invalid, the first return value will be nil, and an
error will be passed.

    func pusherWebhook(res http.ResponseWriter, req *http.Request) {

        body, _ := ioutil.ReadAll(req.Body)
        webhook, err := client.Webhook(req.Header, body)
        if err != nil {
          fmt.Println("Webhook is invalid :(")
        } else {
          fmt.Printf("%+v\n", webhook.Events)
        }

    }
*/
func (c *Client) Webhook(header http.Header, body []byte) (*Webhook, error) {
	for _, token := range header["X-Pusher-Key"] {
		if token == c.Key && checkSignature(header.Get("X-Pusher-Signature"), c.Secret, body) {
			unmarshalledWebhooks, err := unmarshalledWebhook(body)
			if err != nil {
				return unmarshalledWebhooks, err
			}
			decryptedWebhooks, err := decryptEvents(*unmarshalledWebhooks, c.EncryptionMasterKey)
			return decryptedWebhooks, err
		}
	}
	return nil, errors.New("Invalid webhook")
}
