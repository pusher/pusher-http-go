# Pusher Channels HTTP Go Library

[![Build Status](https://travis-ci.org/pusher/pusher-http-go.svg?branch=master)](https://travis-ci.org/pusher/pusher-http-go) [![Coverage Status](https://coveralls.io/repos/pusher/pusher-http-go/badge.svg?branch=master)](https://coveralls.io/r/pusher/pusher-http-go?branch=master) [![GoDoc](https://godoc.org/github.com/pusher/pusher-http-go?status.svg)](https://godoc.org/github.com/pusher/pusher-http-go)

The Golang library for interacting with the Pusher Channels HTTP API.

This package lets you trigger events to your client and query the state of your Pusher channels. When used with a server, you can validate Pusher Channels webhooks and authenticate `private-` or `presence-` channels.

Register for free at <https://pusher.com/channels> and use the application credentials within your app as shown below.

## Supported Platforms

* Go - supports **Go 1.5 or greater**.

## Table of Contents

- [Installation](#installation)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
  - [Additional options](#additional-options)
  - [Google App Engine](#google-app-engine)
- [Usage](#usage)
  - [Triggering events](#triggering-events)
  - [Excluding event recipients](#excluding-event-recipients)
  - [Authenticating Channels](#authenticating-channels)
  - [Application state](#application-state)
  - [Webhook validation](#webhook-validation)
- [Feature Support](#feature-support)
- [Developing the Library](#developing-the-library)
  - [Running the tests](#running-the-tests)
- [License](#license)

## Installation

```sh
$ go get github.com/pusher/pusher-http-go
```

## Getting Started

```go
package main

import (
  "github.com/pusher/pusher-http-go"
)

func main(){
    // instantiate a client
    client := pusher.Client{
        AppID:   "APP_ID",
        Key:     "APP_KEY",
        Secret:  "APP_SECRET",
        Cluster: "APP_CLUSTER",
    }

    data := map[string]string{"message": "hello world"}

    // trigger an event on a channel, along with a data payload
    err := client.Trigger("my-channel", "my_event", data)

    // All trigger methods return an error object, it's worth at least logging this!
    if err != nil {
        panic(err)
    }
}
```

## Configuration

The easiest way to configure the library is by creating a new `Pusher` instance:

```go
client := pusher.Client{
    AppID:   "APP_ID",
    Key:     "APP_KEY",
    Secret:  "APP_SECRET",
    Cluster: "APP_CLUSTER",
}
```

### Additional options

#### Instantiation From URL

```go
client := pusher.ClientFromURL("http://<key>:<secret>@api-<cluster>.pusher.com/apps/app_id")
```

Note: the API URL differs depending on the cluster your app was created in:

```
http://key:secret@api-eu.pusher.com/apps/app_id
http://key:secret@api-ap1.pusher.com/apps/app_id
```

#### Instantiation From Environment Variable

```go
client := pusher.ClientFromEnv("PUSHER_URL")
```

This is particularly relevant if you are using Pusher Channels as a Heroku add-on, which stores credentials in a `"PUSHER_URL"` environment variable.

#### HTTPS

To ensure requests occur over HTTPS, set the `Secure` property of a `pusher.Client` to `true`.

```go
client.Secure = true
```

This is `false` by default.

#### Request Timeouts

If you wish to set a time-limit for each HTTP request, create a `http.Client` instance with your specified `Timeout` field and set it as the Pusher Channels instance's `Client`:

```go
httpClient := &http.Client{Timeout: time.Second * 3}

pusherClient.HTTPClient = httpClient
```

If you do not specifically set a HTTP client, a default one is created with a timeout of 5 seconds.

#### Changing Host

Changing the `pusher.Client`'s `Host` property will make sure requests are sent to your specified host.

```go
client.Host = "foo.bar.com"
```

By default, this is `"api.pusherapp.com"`.

#### Changing the Cluster

Setting the `pusher.Client`'s `Cluster` property will make sure requests are sent to the cluster where you created your app.

*NOTE! If `Host` is set then `Cluster` will be ignored.

```go
client.Cluster = "eu" // in this case requests will be made to api-eu.pusher.com.
```
#### End to End Encryption

This library supports end to end encryption of your private channels. This means that only you and your connected clients will be able to read your messages. Pusher cannot decrypt them. You can enable this feature by following these steps:

1. You should first set up Private channels. This involves [creating an authentication endpoint on your server](https://pusher.com/docs/authenticating_users).

2. Next, Specify your 32 character `EncryptionMasterKey`. This is secret and you should never share this with anyone. Not even Pusher.

```go
client := pusher.Client{
    AppID:              "APP_ID",
    Key:                "APP_KEY",
    Secret:             "APP_SECRET",
    Cluster:            "APP_CLUSTER",
    EncryptionMasterKey "abcdefghijklmnopqrstuvwxyzabcdef",
}
```
3. Channels where you wish to use end to end encryption should be prefixed with `private-encrypted-`.

4. Subscribe to these channels in your client, and you're done! You can verify it is working by checking out the debug console on the https://dashboard.pusher.com/ and seeing the scrambled ciphertext.

**Important note: This will not encrypt messages on channels that are not prefixed by private-encrypted-.**

### Google App Engine

As of version 1.0.0, this library is compatible with Google App Engine's urlfetch library. Simply pass in the HTTP client returned by `urlfetch.Client` to your Pusher Channels initialization struct.

```go
package helloworldapp

import (
    "appengine"
    "appengine/urlfetch"
    "fmt"
    "github.com/pusher/pusher-http-go"
    "net/http"
)

func init() {
    http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    urlfetchClient := urlfetch.Client(c)

    client := pusher.Client{
        AppID:      "APP_ID",
        Key:        "APP_KEY",
        Secret:     "APP_SECRET",
        HTTPClient: urlfetchClient,
    }

    client.Trigger("my-channel", "my_event", map[string]string{"message": "hello world"})

    fmt.Fprint(w, "Hello, world!")
}
```

## Usage

### Triggering events

It is possible to trigger an event on one or more channels. Channel names can contain only characters which are alphanumeric, `_` or `-` and have to be at most 200 characters long. Event name can be at most 200 characters long too.

#### Single channel

##### `func (c *Client) Trigger`

| Argument |Description   |
| :-: | :-: |
| channel `string` | The name of the channel you wish to trigger on. |
| event `string` | The name of the event you wish to trigger |
| data `interface{}` | The payload you wish to send. Must be marshallable into JSON. |

```go
data := map[string]string{"hello": "world"}
client.Trigger("greeting_channel", "say_hello", data)
```

#### Multiple channels

##### `func (c. *Client) TriggerMulti`

| Argument | Description |
| :-: | :-: |
| channels `[]string` | A slice of channel names you wish to send an event on. The maximum length is 10. |
| event `string` | As above. |
| data `interface{}` | As above. |

###### Example

```go
client.TriggerMulti([]string{"a_channel", "another_channel"}, "event", data)
```

#### Excluding event recipients

`func (c *Client) TriggerExclusive` and `func (c *Client) TriggerMultiExclusive` follow the patterns above, except a `socket_id` is given as the last parameter.

These methods allow you to exclude a recipient whose connection has that `socket_id` from receiving the event. You can read more [here](http://pusher.com/docs/duplicates).

##### Examples

**On one channel**:

```go
client.TriggerExclusive("a_channel", "event", data, "123.12")
```

**On multiple channels**:

```go
client.TriggerMultiExclusive([]string{"a_channel", "another_channel"}, "event", data, "123.12")
```

#### Batches

##### `func (c. *Client) TriggerBatch`

| Argument | Description |
| :-: | :-: |
| batch `[]Event` | A list of events to publish |

###### Example

```go
client.TriggerBatch([]pusher.Event{
  { Channel: "a_channel", Name: "event", Data: "hello world", nil },
  { Channel: "a_channel", Name: "event", Data: "hi my name is bob", nil },
})
```

### Authenticating Channels

Application security is very important so Pusher Channels provides a mechanism for authenticating a user’s access to a channel at the point of subscription.

This can be used both to restrict access to private channels, and in the case of presence channels notify subscribers of who else is also subscribed via presence events.

This library provides a mechanism for generating an authentication signature to send back to the client and authorize them.

For more information see our [docs](http://pusher.com/docs/authenticating_users).

#### Private channels

##### `func (c *Client) AuthenticatePrivateChannel`

| Argument | Description |
| :-: | :-: |
| params `[]byte` | The request body sent by the client |

| Return Value | Description |
| :-: | :-: |
| response `[]byte` | The response to send back to the client, carrying an authentication signature |
| err `error` | Any errors generated |

###### Example

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
    http.HandleFunc("/pusher/auth", pusherAuth)
    http.ListenAndServe(":5000", nil)
}
```

###### Example (JSONP)

```go
func pusherJsonpAuth(res http.ResponseWriter, req *http.Request) {
    var (
        callback, params string
    )

    {
        q := r.URL.Query()
        callback = q.Get("callback")
        if callback == "" {
            panic("callback missing")
        }
        q.Del("callback")
        params = []byte(q.Encode())
    }

    response, err := client.AuthenticatePrivateChannel(params)
    if err != nil {
        panic(err)
    }

    res.Header().Set("Content-Type", "application/javascript; charset=utf-8")
    fmt.Fprintf(res, "%s(%s);", callback, string(response))
}

func main() {
    http.HandleFunc("/pusher/auth", pusherJsonpAuth)
    http.ListenAndServe(":5000", nil)
}
```

#### Authenticating presence channels

Using presence channels is similar to private channels, but in order to identify a user, clients are sent a user_id and, optionally, custom data.

##### `func (c *Client) AuthenticatePresenceChannel`

| Argument | Description |
| :-: | :-: |
| params `[]byte` | The request body sent by the client |
| member `pusher.MemberData` | A struct representing what to assign to a channel member, consisting of a `UserID` and any custom `UserInfo`. See below |

###### Custom Types

**pusher.MemberData**

```go
type MemberData struct {
    UserID   string
    UserInfo map[string]string
}
```

###### Example

```go
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
```

### Application state

This library allows you to query our API to retrieve information about your application's channels, their individual properties, and, for presence-channels, the users currently subscribed to them.

#### Get the list of channels in an application

##### `func (c *Client) Channels`

| Argument | Description |
| :-: | :-: |
| additionalQueries `map[string]string` | A map with query options. A key with `"filter_by_prefix"` will filter the returned channels. To get number of users subscribed to a presence-channel, specify an `"info"` key with value `"user_count"`. <br><br>Pass in `nil` if you do not wish to specify any query attributes. |

| Return Value | Description |
| :-: | :-: |
| channels `*pusher.ChannelsList` | A struct representing the list of channels. See below. |
| err `error` | Any errors encountered|

###### Custom Types

**pusher.ChannelsList**

```go
type ChannelsList struct {
    Channels map[string]ChannelListItem
}
```

**pusher.ChannelsListItem**

```go
type ChannelListItem struct {
    UserCount int
}
```

###### Example

```go
channelsParams := map[string]string{
    "filter_by_prefix": "presence-",
    "info":             "user_count",
}

channels, err := client.Channels(channelsParams)

// channels => &{Channels:map[presence-chatroom:{UserCount:4} presence-notifications:{UserCount:31}]}
```

#### Get the state of a single channel

##### `func (c *Client) Channel`

| Argument | Description |
| :-: | :-: |
| name `string` | The name of the channel |
| additionalQueries `map[string]string` | A map with query options. An `"info"` key can have comma-separated values of `"user_count"`, for presence-channels, and `"subscription_count"`, for all-channels. To use the `"subscription_count"` value, first check the "Enable subscription counting" checkbox in your App Settings on [your Pusher Channels dashboard](https://dashboard.pusher.com).<br><br>Pass in `nil` if you do not wish to specify any query attributes. |

| Return Value | Description |
| :-: | :-: |
| channel `*pusher.Channel` | A struct representing a channel. See below. |
| err `error` | Any errors encountered |

###### Custom Types

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

// channel => &{Name:presence-chatroom Occupied:true UserCount:42 SubscriptionCount:42}
```

#### Get a list of users in a presence channel

##### `func (c *Client) GetChannelUsers`

| Argument | Description |
| :-: | :-: |
| name `string` | The channel name |

| Return Value | Description |
| :-: | :-: |
| users `*pusher.Users` | A struct representing a list of the users subscribed to the presence-channel. See below |
| err `error` | Any errors encountered. |

###### Custom Types

**pusher.Users**

```go
type Users struct {
    List []User
}
```

**pusher.User**

```go
type User struct {
    ID string
}
```

###### Example

```go
users, err := client.GetChannelUsers("presence-chatroom")

// users => &{List:[{ID:13} {ID:90}]}
```

### Webhook validation

On your [dashboard](http://app.pusher.com), you can set up webhooks to POST a payload to your server after certain events. Such events include channels being occupied or vacated, members being added or removed in presence-channels, or after client-originated events. For more information see <https://pusher.com/docs/webhooks>.

This library provides a mechanism for checking that these POST requests are indeed from Pusher, by checking the token and authentication signature in the header of the request.

##### `func (c *Client) Webhook`

| Argument | Description |
| :-: | :-: |
| header `http.Header` | The header of the request to verify |
| body `[]byte` | The body of the request |

| Return Value | Description |
| :-: | :-: |
| webhook `*pusher.Webhook` | If the webhook is valid, this method will return a representation of that webhook that includes its timestamp and associated events. If invalid, this value will be `nil`. |
| err `error` | If the webhook is invalid, an error value will be passed. |

###### Custom Types

**pusher.Webhook**

```go
type Webhook struct {
    TimeMs int
    Events []WebhookEvent
}
```

**pusher.WebhookEvent**

```go
type WebhookEvent struct {
    Name     string
    Channel  string
    Event    string
    Data     string
    SocketID string
}
```

###### Example

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

## Feature Support

Feature                                    | Supported
-------------------------------------------| :-------:
Trigger event on single channel            | *&#10004;*
Trigger event on multiple channels         | *&#10004;*
Trigger events in batches                  | *&#10004;*
Excluding recipients from events           | *&#10004;*
Authenticating private channels            | *&#10004;*
Authenticating presence channels           | *&#10004;*
Get the list of channels in an application | *&#10004;*
Get the state of a single channel          | *&#10004;*
Get a list of users in a presence channel  | *&#10004;*
WebHook validation                         | *&#10004;*
Heroku add-on support                      | *&#10004;*
Debugging & Logging                        | *&#10004;*
Cluster configuration                      | *&#10004;*
Timeouts                                   | *&#10004;*
HTTPS                                      | *&#10004;*
HTTP Proxy configuration                   | *&#10008;*
HTTP KeepAlive                             | *&#10008;*

## Helper Functionality

These are helpers that have been implemented to to ensure interactions with the HTTP API only occur if they will not be rejected e.g. [channel naming conventions](https://pusher.com/docs/channels/using_channels/channels#channel-naming-conventions).

Helper Functionality                      | Supported
----------------------------------------- | :-------:
Channel name validation                   | &#10004;
Limit to 10 channels per trigger          | &#10004;
Limit event name length to 200 chars      | &#10004;

## Developing the Library

Feel more than free to fork this repo, improve it in any way you'd prefer, and send us a pull request :)

### Running the tests

Simply type:

```sh
$ go test
```

## License

This code is free to use under the terms of the MIT license.
