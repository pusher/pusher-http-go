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

const (
	libraryVersion = "5.1.1"
	libraryName    = "pusher-http-go"
)

/*
Client to the HTTP API of Pusher.

There easiest way to configure the library is by creating a `Pusher` instance:

	client := pusher.Client{
		AppID: "your_app_id",
		Key: "your_app_key",
		Secret: "your_app_secret",
	}

To ensure requests occur over HTTPS, set the `Secure` property of a
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
	AppID                        string
	Key                          string
	Secret                       string
	Host                         string // host or host:port pair
	Secure                       bool   // true for HTTPS
	Cluster                      string
	HTTPClient                   *http.Client
	EncryptionMasterKey          string  // deprecated
	EncryptionMasterKeyBase64    string  // for E2E
	OverrideMaxMessagePayloadKB  int     // set the agreed Pusher message limit increase
	validatedEncryptionMasterKey *[]byte // parsed key for use
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
	_, err := c.validateChannelsAndTrigger([]string{channel}, eventName, data, TriggerParams{})
	return err
}

/*
ChannelsParams are any parameters than can be sent with a
TriggerWithParams or TriggerMultiWithParams requests.
*/
type TriggerParams struct {
	// SocketID excludes a recipient whose connection has the `socket_id`
	// specified here. You can read more here:
	// http://pusher.com/docs/duplicates.
	SocketID *string
	// Info is comma-separated vales of `"user_count"`, for
	// presence-channels, and `"subscription_count"`, for all-channels.
	// Note that the subscription count is not allowed by default. Please
	// contact us at http://support.pusher.com if you wish to enable this.
	// Pass in `nil` if you do not wish to specify any query attributes.
	// This is part of an [experimental feature](https://pusher.com/docs/lab#experimental-program).
	Info *string
}

func (params TriggerParams) toMap() map[string]string {
	m := make(map[string]string)
	if params.SocketID != nil {
		m["socket_id"] = *params.SocketID
	}
	if params.Info != nil {
		m["info"] = *params.Info
	}
	return m
}

/*
TriggerWithParams is the same as `client.Trigger`, except it allows additional
parameters to be passed in. See:
https://pusher.com/docs/channels/library_auth_reference/rest-api#request
for a complete list.

	data := map[string]string{"hello": "world"}
	socketID := "1234.12"
	attributes := "user_count"
	params := pusher.TriggerParams{SocketID: &socketID, Info: &attributes}
	channels, err := client.Trigger("greeting_channel", "say_hello", data, params)

	//channels=> &{Channels:map[presence-chatroom:{UserCount:4} presence-notifications:{UserCount:31}]}
*/
func (c *Client) TriggerWithParams(
	channel string,
	eventName string,
	data interface{},
	params TriggerParams,
) (*TriggerChannelsList, error) {
	return c.validateChannelsAndTrigger([]string{channel}, eventName, data, params)
}

/*
TriggerMulti is the same as `client.Trigger`, except one passes in a slice of
`channels` as the first parameter. The maximum length of channels is 100.

	client.TriggerMulti([]string{"a_channel", "another_channel"}, "event", data)
*/
func (c *Client) TriggerMulti(channels []string, eventName string, data interface{}) error {
	_, err := c.validateChannelsAndTrigger(channels, eventName, data, TriggerParams{})
	return err
}

/*
TriggerMultiWithParams is the same as `client.TriggerMulti`, except it
allows additional parameters to be specified in the same way as
`client.TriggerWithParams`.
*/
func (c *Client) TriggerMultiWithParams(
	channels []string,
	eventName string,
	data interface{},
	params TriggerParams,
) (*TriggerChannelsList, error) {
	return c.validateChannelsAndTrigger(channels, eventName, data, params)
}

/*
TriggerExclusive triggers an event excluding a recipient whose connection has
the `socket_id` you specify here from receiving the event.
You can read more here: http://pusher.com/docs/duplicates.

	client.TriggerExclusive("a_channel", "event", data, "123.12")

Deprecated: use TriggerWithParams instead.
*/
func (c *Client) TriggerExclusive(channel string, eventName string, data interface{}, socketID string) error {
	params := TriggerParams{SocketID: &socketID}
	_, err := c.validateChannelsAndTrigger([]string{channel}, eventName, data, params)
	return err
}

/*
TriggerMultiExclusive triggers an event to multiple channels excluding a
recipient whose connection has the `socket_id` you specify here from receiving
the event on any of the channels.

	client.TriggerMultiExclusive([]string{"a_channel", "another_channel"}, "event", data, "123.12")

Deprecated: use TriggerMultiWithParams instead.
*/
func (c *Client) TriggerMultiExclusive(channels []string, eventName string, data interface{}, socketID string) error {
	params := TriggerParams{SocketID: &socketID}
	_, err := c.validateChannelsAndTrigger(channels, eventName, data, params)
	return err
}

/*
SendToUser triggers an event to a specific user.
Pass in the user id, the event's name, and a data payload. The data payload must
be marshallable into JSON.

	data := map[string]string{"hello": "world"}
	client.SendToUser("user123", "say_hello", data)
*/
func (c *Client) SendToUser(userId string, eventName string, data interface{}) error {
	if !validUserId(userId) {
		return fmt.Errorf("User id '%s' is invalid", userId)
	}
	_, err := c.trigger([]string{"#server-to-user-" + userId}, eventName, data, TriggerParams{})
	return err
}

func (c *Client) validateChannelsAndTrigger(channels []string, eventName string, data interface{}, params TriggerParams) (*TriggerChannelsList, error) {
	if len(channels) > maxTriggerableChannels {
		return nil, fmt.Errorf("You cannot trigger on more than %d channels at once", maxTriggerableChannels)
	}
	if !channelsAreValid(channels) {
		return nil, errors.New("At least one of your channels' names are invalid")
	}
	return c.trigger(channels, eventName, data, params)
}

func (c *Client) trigger(channels []string, eventName string, data interface{}, params TriggerParams) (*TriggerChannelsList, error) {
	hasEncryptedChannel := false
	for _, channel := range channels {
		if isEncryptedChannel(channel) {
			hasEncryptedChannel = true
		}
	}
	if hasEncryptedChannel && len(channels) > 1 {
		// For rationale, see limitations of end-to-end encryption in the README
		return nil, errors.New("You cannot trigger to multiple channels when using encrypted channels")
	}
	masterKey, keyErr := c.encryptionMasterKey()
	if hasEncryptedChannel && keyErr != nil {
		return nil, keyErr
	}

	if err := validateSocketID(params.SocketID); err != nil {
		return nil, err
	}

	payload, err := encodeTriggerBody(channels, eventName, data, params.toMap(), masterKey, c.OverrideMaxMessagePayloadKB)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/apps/%s/events", c.AppID)
	triggerURL, err := createRequestURL("POST", c.Host, path, c.Key, c.Secret, authTimestamp(), c.Secure, payload, nil, c.Cluster)
	if err != nil {
		return nil, err
	}
	response, err := c.request("POST", triggerURL, payload)
	if err != nil {
		return nil, err
	}

	return unmarshalledTriggerChannelsList(response)
}

/*
Event stores all the data for one Event that can be triggered.
*/
type Event struct {
	Channel  string
	Name     string
	Data     interface{}
	SocketID *string
	// Info is part of an [experimental feature](https://pusher.com/docs/lab#experimental-program).
	Info *string
}

/*
TriggerBatch triggers multiple events on multiple channels in a single call:

	info := "subscription_count"
	socketID := "1234.12"
	client.TriggerBatch([]pusher.Event{
		{ Channel: "donut-1", Name: "ev1", Data: "d1", SocketID: socketID, Info: &info },
		{ Channel: "private-encrypted-secretdonut", Name: "ev2", Data: "d2", SocketID: socketID, Info: &info },
	})
*/
func (c *Client) TriggerBatch(batch []Event) (*TriggerBatchChannelsList, error) {
	hasEncryptedChannel := false
	// validate every channel name and every sockedID (if present) in batch
	for _, event := range batch {
		if !validChannel(event.Channel) {
			return nil, fmt.Errorf("The channel named %s has a non-valid name", event.Channel)
		}
		if err := validateSocketID(event.SocketID); err != nil {
			return nil, err
		}
		if isEncryptedChannel(event.Channel) {
			hasEncryptedChannel = true
		}
	}
	masterKey, keyErr := c.encryptionMasterKey()
	if hasEncryptedChannel && keyErr != nil {
		return nil, keyErr
	}

	payload, err := encodeTriggerBatchBody(batch, masterKey, c.OverrideMaxMessagePayloadKB)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/apps/%s/batch_events", c.AppID)
	triggerURL, err := createRequestURL("POST", c.Host, path, c.Key, c.Secret, authTimestamp(), c.Secure, payload, nil, c.Cluster)
	if err != nil {
		return nil, err
	}
	response, err := c.request("POST", triggerURL, payload)
	if err != nil {
		return nil, err
	}

	return unmarshalledTriggerBatchChannelsList(response)
}

/*
ChannelsParams are any parameters than can be sent with a Channels request.
*/
type ChannelsParams struct {
	// FilterByPrefix will filter the returned channels.
	FilterByPrefix *string
	// Info should be specified with a value of "user_count" to get number
	// of users subscribed to a presence-channel. Pass in `nil` if you do
	// not wish to specify any query attributes.
	Info *string
}

func (params ChannelsParams) toMap() map[string]string {
	m := make(map[string]string)
	if params.FilterByPrefix != nil {
		m["filter_by_prefix"] = *params.FilterByPrefix
	}
	if params.Info != nil {
		m["info"] = *params.Info
	}
	return m
}

/*
Channels returns a list of all the channels in an application.

	prefixFilter := "presence-"
	attributes := "user_count"
	params := pusher.ChannelsParams{FilterByPrefix: &prefixFilter, Info: &attributes}
	channels, err := client.Channels(params)

	//channels=> &{Channels:map[presence-chatroom:{UserCount:4} presence-notifications:{UserCount:31}  ]}
*/
func (c *Client) Channels(params ChannelsParams) (*ChannelsList, error) {
	path := fmt.Sprintf("/apps/%s/channels", c.AppID)
	u, err := createRequestURL("GET", c.Host, path, c.Key, c.Secret, authTimestamp(), c.Secure, nil, params.toMap(), c.Cluster)
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
ChannelParams are any parameters than can be sent with a Channel request.
*/
type ChannelParams struct {
	// Info is comma-separated vales of `"user_count"`, for
	// presence-channels, and `"subscription_count"`, for all-channels.
	// Note that the subscription count is not allowed by default. Please
	// contact us at http://support.pusher.com if you wish to enable this.
	// Pass in `nil` if you do not wish to specify any query attributes.
	Info *string
}

func (params ChannelParams) toMap() map[string]string {
	m := make(map[string]string)
	if params.Info != nil {
		m["info"] = *params.Info
	}
	return m
}

/*
Channel allows you to get the state of a single channel.

	attributes := "user_count,subscription_count"
	params := pusher.ChannelParams{Info: &attributes}
	channel, err := client.Channel("presence-chatroom", params)

	//channel=> &{Name:presence-chatroom Occupied:true UserCount:42 SubscriptionCount:42}
*/
func (c *Client) Channel(name string, params ChannelParams) (*Channel, error) {
	path := fmt.Sprintf("/apps/%s/channels/%s", c.AppID, name)
	u, err := createRequestURL("GET", c.Host, path, c.Key, c.Secret, authTimestamp(), c.Secure, nil, params.toMap(), c.Cluster)
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
AuthenticateUser allows you to authenticate a user's connection.
It returns an authentication signature to send back to the client
and authenticate them. In order to identify a user, this method acceps a map containing
arbitrary user data. It must contain at least an id field with the user's id as a string.

For more information see our docs: http://pusher.com/docs/authenticating_users.

This is an example of authenticating a user, using the built-in
Golang HTTP library to start a server.

In order to authenticate a client, one must read the response into type `[]byte`
and pass it in. This will return a signature in the form of a `[]byte` for you
to send back to the client.

	func pusherUserAuth(res http.ResponseWriter, req *http.Request) {

		params, _ := ioutil.ReadAll(req.Body)
		userData := map[string]interface{} { "id": "1234", "twitter": "jamiepatel" }
		response, err := client.AuthenticateUser(params, userData)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(res, string(response))
	}

	func main() {
		http.HandleFunc("/pusher/user-auth", pusherUserAuth)
		http.ListenAndServe(":5000", nil)
	}
*/
func (c *Client) AuthenticateUser(params []byte, userData map[string]interface{}) (response []byte, err error) {
	socketID, err := parseUserAuthenticationRequestParams(params)
	if err != nil {
		return
	}

	if err = validateSocketID(&socketID); err != nil {
		return
	}

	if err = validateUserData(userData); err != nil {
		return
	}

	var jsonUserData string
	if jsonUserData, err = jsonMarshalToString(userData); err != nil {
		return
	}
	stringToSign := strings.Join([]string{socketID, "user", jsonUserData}, "::")

	_response := createAuthMap(c.Key, c.Secret, stringToSign, "")
	_response["user_data"] = jsonUserData

	response, err = json.Marshal(_response)
	return
}

/*
AuthorizePrivateChannel allows you to authorize a users subscription to a
private channel. It returns an authorization signature to send back to the client
and authorize them.

For more information see our docs: http://pusher.com/docs/authorizing_users.

This is an example of authorizing a private-channel, using the built-in
Golang HTTP library to start a server.

In order to authorize a client, one must read the response into type `[]byte`
and pass it in. This will return a signature in the form of a `[]byte` for you
to send back to the client.

	func pusherAuth(res http.ResponseWriter, req *http.Request) {

		params, _ := ioutil.ReadAll(req.Body)
		response, err := client.AuthorizePrivateChannel(params)
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
func (c *Client) AuthorizePrivateChannel(params []byte) (response []byte, err error) {
	return c.authorizeChannel(params, nil)
}

/*
AuthenticatePrivateChannel allows you to authorize a users subscription to a
private channel. It returns an authorization signature to send back to the client
and authorize them.

Deprecated: use AuthorizePrivateChannel instead.
*/
func (c *Client) AuthenticatePrivateChannel(params []byte) (response []byte, err error) {
	return c.authorizeChannel(params, nil)
}

/*
AuthorizePresenceChannel allows you to authorize a users subscription to a
presence channel. It returns an authorization signature to send back to the client
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

	response, err := client.AuthorizePresenceChannel(params, presenceData)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(res, response)
*/
func (c *Client) AuthorizePresenceChannel(params []byte, member MemberData) (response []byte, err error) {
	return c.authorizeChannel(params, &member)
}

/*
AuthenticatePresenceChannel allows you to authorize a users subscription to a
presence channel. It returns an authorization signature to send back to the client
and authorize them. In order to identify a user, clients are sent a user_id and,
optionally, custom data.

Deprecated: use AuthorizePresenceChannel instead.
*/
func (c *Client) AuthenticatePresenceChannel(params []byte, member MemberData) (response []byte, err error) {
	return c.authorizeChannel(params, &member)
}

func (c *Client) authorizeChannel(params []byte, member *MemberData) (response []byte, err error) {
	channelName, socketID, err := parseChannelAuthorizationRequestParams(params)
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
		masterKey, err := c.encryptionMasterKey()
		if err != nil {
			return nil, err
		}
		sharedSecret := generateSharedSecret(channelName, masterKey)
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
				return nil, err
			}

			hasEncryptedChannel := false
			for _, event := range unmarshalledWebhooks.Events {
				if isEncryptedChannel(event.Channel) {
					hasEncryptedChannel = true
				}
			}
			masterKey, keyErr := c.encryptionMasterKey()
			if hasEncryptedChannel && keyErr != nil {
				return nil, keyErr
			}

			return decryptEvents(*unmarshalledWebhooks, masterKey)
		}
	}
	return nil, errors.New("Invalid webhook")
}

func (c *Client) encryptionMasterKey() ([]byte, error) {
	if c.validatedEncryptionMasterKey != nil {
		return *(c.validatedEncryptionMasterKey), nil
	}

	if c.EncryptionMasterKey != "" && c.EncryptionMasterKeyBase64 != "" {
		return nil, errors.New("Do not specify both EncryptionMasterKey and EncryptionMasterKeyBase64. EncryptionMasterKey is deprecated, specify only EncryptionMasterKeyBase64")
	}

	if c.EncryptionMasterKey != "" {
		if len(c.EncryptionMasterKey) != 32 {
			return nil, errors.New("EncryptionMasterKey must be 32 bytes. It is also deprecated, use EncryptionMasterKeyBase64")
		}

		keyBytes := []byte(c.EncryptionMasterKey)
		c.validatedEncryptionMasterKey = &keyBytes
		return keyBytes, nil
	}

	if c.EncryptionMasterKeyBase64 != "" {
		keyBytes, err := base64.StdEncoding.DecodeString(c.EncryptionMasterKeyBase64)
		if err != nil {
			return nil, errors.New("EncryptionMasterKeyBase64 must be valid base64")
		}
		if len(keyBytes) != 32 {
			return nil, errors.New("EncryptionMasterKeyBase64 must encode 32 bytes")
		}

		c.validatedEncryptionMasterKey = &keyBytes
		return keyBytes, nil
	}

	return nil, errors.New("No master encryption key supplied")
}
