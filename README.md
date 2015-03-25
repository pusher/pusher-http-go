# Pusher HTTP Go Library

The Golang library for interacting with the Pusher HTTP API.

This package lets you trigger events to your client and query the state of your Pusher channels. When used with a server, you can validate Pusher webhooks and authenticate private- or presence-channels.

In order to use this library, you need to have an account on <http://pusher.com>. After registering, you will need the application credentials for your app.


## Installation

```
$ go get github.com/pusher/pusher-http-go
```

## Getting Started

```go
// instantiate a client
client := pusher.Client{
  AppId: "your_app_id",
  Key: "your_app_key",
  Secret: "your_app_secret",
}

data := map[string]string{"message": "hello world"}

// trigger an event on a channel, along with a data payload
client.Trigger("test_channel", "event", data)
```

## Configuration

There easiest way to configure the library is by creating a new `Pusher` instance:

```go
client := pusher.Client{
  AppId: "your_app_id",
  Key: "your_app_key",
  Secret: "your_app_secret",
}
```

### Additional options

#### From URL

```go
client := pusher.ClientFromURL("http://key:secret@api.pusherapp.com/apps/app_id")
```

#### From Environment Variable

```go
client := pusher.ClientFromEnv("PUSHER_URL")
```

This is particularly relevant if you are using Pusher as a Heroku add-on, which stores credentials in a `"PUSHER_URL"` environment variable.

## Usage

### Triggering events

It is possible to trigger an event on one or more channels. Channel names can contain only characters which are alphanumeric, `_` or `-`` and have to be at most 200 characters long. Event name can be at most 200 characters long too.


#### Single channel

#####`func (c *Client) Trigger`

|Argument   |Description   |
|:-:|:-:|
|channel `string`   |The name of the channel you wish to trigger on.   |
|event `string` | The name of the event you wish to trigger |
|data `interface{}` | The payload you wish to send. Must be marshallable into JSON. |

```go
data := map[string]string{"hello": "world"}
client.Trigger("my_numbers", "numbers_for_all", data)
```

#### Multiple channels

#####`func (c. *Client) TriggerMulti`

|Argument | Description |
|:-:|:-:|
|channels `[]string`| A slice of channel names you wish to send an event on. The maximum length is 10.|
|event `string` | As above.|
|data `interface{}` |As above.|

######Example

```go
client.TriggerMulti([]string{"a_channel", "another_channel"}, "event", data)
```

### Excluding event recipients

`func (c *Client) TriggerExclusive` and `func (c *Client) TriggerMultiExclusive` follow the patterns above, except a `socket_id` is given as the last parameter.

These methods allow you to exclude a recipient whose connection has that `socket_id` from receiving the event. You can read more [here](http://pusher.com/docs/duplicates).

######Examples

**On one channel**:

```go
client.TriggerExclusive("a_channel", "event", data, "123.12")
```

**On multiple channels**:

```go
client.TriggerMultiExclusive([]string{"a_channel", "another_channel"}, "event", data, "123.12")
```

### Authenticating Channels

Application security is very important so Pusher provides a mechanism for authenticating a user’s access to a channel at the point of subscription.

This can be used both to restrict access to private channels, and in the case of presence channels notify subscribers of who else is also subscribed via presence events.

This library provides a mechanism for generating an authentication signature to send back to the client and authorize them.

For more information see our [docs](http://pusher.com/docs/authenticating_users).

#### Private channels


##### `func (c *Client) AuthenticatePrivateChannel`

|Argument|Description|
|:-:|:-:|
|params `[]byte`| The request body sent by the client|

|Return Values|Description|
|:-:|:-:|
|response `[]byte` | The response to send back to the client, carrying an authentication signature |
|err `error` | Any errors generated |

######Example

```go
func pusherAuth(res http.ResponseWriter, req *http.Request) {

	params, _ := ioutil.ReadAll(req.Body)
	response, err := client.AuthenticatePrivateChannel(params)
	
	if err != nil {
		panic(err)
	}
	
	fmt.Fprintf(res, string(response))

}

func main() {
	http.HandleFunc("/pusher/auth", pusher_auth)
	http.ListenAndServe(":5000", nil)
}
```

#### Authenticating presence channels

Using presence channels is similar to private channels, but in order to identify a user, clients are sent a user_id and, optionally, custom data.

##### `func (c *Client) AuthenticatePresenceChannel`

|Argument|Description|
|:-:|:-:|
|params `[]byte`| The request body sent by the client |
|member `pusher.MemberData`| A struct representing what to assign to a channel member, consisting of a `UserId` and any custom `UserInfo`. See below |

###### Custom Types

**pusher.MemberData**

```go
type MemberData struct {
    UserId   string            `json:"user_id"`
    UserInfo map[string]string `json:"user_info",omitempty`
}
```

###### Example

```go
params, _ := ioutil.ReadAll(req.Body)

presenceData := pusher.MemberData{
	UserId: "1",
	UserInfo: map[string]string{
		"twitter": "jamiepatel",
	},
}

response, err := client.AuthenticatePresenceChannel(params, presenceData)

if err != nil {
	panic(err)
}

fmt.Fprintf(res, response)
```

### Application state

This library allows you to query our API to retrieve information about your application's channels, their individual properties, and, for presence-channels, the users currently subscribed to them.

#### Get the list of channels in an application

##### `func (c *Client) Channels`

|Argument|Description|
|:-:|:-:|
|additionalQueries `map[string]string`| A map with query options. A key with `"filter_by_prefix"` will filter the returned channels. To get number of users subscribed to a presence-channel, specify an `"info"` key with value `"user_count"`. <br><br>Pass in `nil` if you do not wish to specify any query attributes.  |

|Return Values|Description|
|:-:|:-:|
|channels `*pusher.ChannelsList`|A struct representing the list of channels. See below. |
|err `error`|Any errors encountered|

###### Custom Types

**pusher.ChannelsList**

```go
type ChannelsList struct {
    Channels map[string]ChannelListItem `json:"channels"`
}
```

**pusher.ChannelsListItem**

```go
type ChannelListItem struct {
    UserCount int `json:"user_count"`
}
```
######Example

```go
channelsParams := map[string]string{
    "filter_by_prefix": "presence-",
    "info":             "user_count",
}

channels, err := client.Channels(channelsParams)

// => &{Channels:map[presence-chatroom:{UserCount:4} presence-notifications:{UserCount:31}  ]}
```

#### Get the state of a single channel

##### `func (c *Client) Channel`

|Argument|Description|
|:-:|:-:|
|name `string`| The name of the channel|
|additionalQueries `map[string]string` |A map with query options. An `"info"` key can have comma-separated vales of `"user_count"`, for presence-channels, and `"subscription_count"`, for all-channels. Note that the subscription count is not allowed by default. Please [contact us](http://support.pusher.com) if you wish to enable this.<br><br>Pass in `nil` if you do not wish to specify any query attributes.|

|Return Values|Description|
|:-:|:-:|
|channel `*pusher.Channel` |A struct representing a channel. See below. |

######Custom Types

**pusher.Channel**

```go
type Channel struct {
    Name              string
    Occupied          bool
    UserCount         int  
    SubscriptionCount int 
}
```

###### Example
```go
channelParams := map[string]string{
	"info": "user_count,subscription_count",
}

channel, err := client.Channel("presence-chatroom", channelParams)

//=> &{Name:presence-chatroom Occupied:true UserCount:42 SubscriptionCount:42}
```

#### Get a list of users in a presence channel

*Any additional information specific to the library*

**{Example}:**

```go
users, err := client.GetChannelUsers("presence-chatroom")
```

### WebHook validation

*Not all libraries presently offer this functionality. But if they do...*

The library provides a simple helper for WebHooks.

*Any additional information specific to the library*

For more information see <https://pusher.com/docs/webhooks>.

**{Example}:**

```go
func pusherWebhook(res http.ResponseWriter, req *http.Request) {

	body, _ := ioutil.ReadAll(req.Body)
	webhook, err := client.Webhook(req.Header, body)
  if err != nil {
      fmt.Println("Webhook is invalid :(")
  } else {
      fmt.Printf("%+v\n", webhook.Events)
  }

}
```

### Debugging & Logging

*Information on how to debug the library and providing logging information. We've found that this is very useful during the development process*

For additional information on debugging and logging please see <https://pusher.com/docs/debugging>.

### Feature Support

*Provide information regarding the features that the library supports. What it does and what it doesn't. This section can also form a table of contents to the information within the README*


Feature                                    | Supported
-------------------------------------------| :-------:
Trigger event on single channel            | *&#10004;*
Trigger event on multiple channels         | *&#10004;*
Excluding recipients from events           | *&#10004;*
Authenticating private channels            | *&#10004;*
Authenticating presence channels           | *&#10004;*
Get the list of channels in an application | *&#10004;*
Get the state of a single channel          | *&#10004;*
Get a list of users in a presence channel  | *&#10004;*
WebHook validation                         | *&#10004;*
Heroku add-on support						   | *&#10004;*
Debugging & Logging                        | *&#10004;*
Cluster configuration                      | *&#10004;*
Timeouts                                   | *&#10004;*
HTTPS                                      | *&#10004;*
HTTP Proxy configuration                   | *&#10008;*
HTTP KeepAlive                             | *&#10008;*


### Helper Functionality

*Libraries can also offer additional helper functionality to ensure interactions with the HTTP API only occur if they will not be rejected e.g. [channel naming conventions][channel-names].*

Helper Functionality                     | Supported
-----------------------------------------| :-------:
[Channel name validation][channel-names] | &#10004;
Limit to 10 channels per trigger         | &#10004;
Limit event name length to 200 chars     | &#10004;

## Developing the Library

*A section providing information for developers who wish to develop the library*

### Testing

*Any information specific to the library*

### Running tests

    $ go test

### Deploy to Distribution Mechanism

*Any additional information specific to the library*

## Credits

*It's always nice to give credit to those who inspired the work or who have contributed*

## License

This code is free to use under the terms of the MIT license.

[channel-names]: https://pusher.com/docs/client_api_guide/client_channels#naming-channels
