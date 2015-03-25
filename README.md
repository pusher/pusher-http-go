# Pusher HTTP Go library

*Badges e.g. Travis Build, Dependency Status*

The Golang library for interacting with the Pusher HTTP API.

*Brief overview of what the library offers providing more detail than simply interacting with the HTTP API*

In order to use this library, you need to have an account on <http://pusher.com>. After registering, you will need the application credentials for your app.

<!-- ## Feature Support

*Provide information regarding the features that the library supports. What it does and what it doesn't. This section can also form a table of contents to the information within the README*

Feature                                    | Supported
-------------------------------------------| :-------:
Trigger event on single channel            | *&#10004; or &#10008;*
Trigger event on multiple channels         | *&#10004; or &#10008;*
Excluding recipients from events           | *&#10004; or &#10008;*
Authenticating private channels            | *&#10004; or &#10008;*
Authenticating presence channels           | *&#10004; or &#10008;*
Get the list of channels in an application | *&#10004; or &#10008;*
Get the state of a single channel          | *&#10004; or &#10008;*
Get a list of users in a presence channel  | *&#10004; or &#10008;*
WebHook validation                         | *&#10004; or &#10008;*
Debugging & Logging                        | *&#10004; or &#10008;*
HTTPS                                      | *&#10004; or &#10008;*
HTTP Proxy configuration                   | *&#10004; or &#10008;*
Cluster configuration                      | *&#10004; or &#10008;*

Libraries can also offer additional helper functionality to ensure interactions with the HTTP API only occur if they will not be rejected e.g. [channel naming conventions][channel-names]. For information on the helper functionality that this library supports please see the **Helper Functionality** section. -->

## Installation

```
$ go get github.com/pusher/pusher-http-go
```

## Getting Started

```go
client := pusher.Client{
  AppId: "your_app_id",
  Key: "your_app_key",
  Secret: "your_app_secret",
}

data := map[string]string{"message": "hello world"}

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

*Notes on triggering an event on a single channel*

**{Example}:**

```go
data := map[string]string{"hello": "world"}
client.Trigger("my_numbers", "numbers_for_all", data)
```

#### Multiple channels

*Notes on triggering an event on multiple channels*

**Example**

```go
client.TriggerMulti([]string{"a_channel", "another_channel"}, "event", data)
```

### Excluding event recipients

*Notes on triggering an event and identifying a socket_id that will not receive the event*

**{Example}:**

```go
client.TriggerExclusive("a_channel", "event", data, "123.12")
```

```go
client.TriggerMultiExclusive([]string{"a_channel", "another_channel"}, "event", data, "123.12")
```

### Authenticating private channels

To authorize your users to access private channels on Pusher *...*

For more information see: <http://pusher.com/docs/authenticating_users>

**{Example}:**

```go
func pusherAuth(res http.ResponseWriter, req *http.Request) {

	params, _ := ioutil.ReadAll(req.Body)
	response := client.AuthenticatePrivateChannel(params)
	fmt.Fprintf(res, response)

}

func main() {
	http.HandleFunc("/pusher/auth", pusher_auth)
	http.ListenAndServe(":5000", nil)
}
```

### Authenticating presence channels

Using presence channels is similar to private channels, but you can specify extra data to identify that particular user.

*Any additional information specific to the library*

For more information see: <http://pusher.com/docs/authenticating_users>

**{Example}:**

```go
params, _ := ioutil.ReadAll(req.Body)
presenceData := pusher.MemberData{UserId: "1", UserInfo: map[string]string{"twitter": "@pusher"}}
response := client.AuthenticatePresenceChannel(params, presence_data)
fmt.Fprintf(res, response)
```

### Application state

It's possible to query the state of the application.

*Any additional information specific to the library*

**{Example}:**

```js
pusher.get({ path: path, params: params }, callback);
```

#### Get the list of channels in an application

*Any additional information specific to the library*

**Example**:

```go
channels, err := client.Channels(channelsParams)
```

#### Get the state of a single channel

*Any additional information specific to the library*

**{Example}:**

```go
channel, err := client.Channel("presence-chatroom")
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

### Helper Functionality

*Libraries can also offer additional helper functionality to ensure interactions with the HTTP API only occur if they will not be rejected e.g. [channel naming conventions][channel-names].*

Helper Functionality                     | Supported
-----------------------------------------| :-------:
[Channel name validation][channel-names] | *&#10004; or &#10008;*
Limit to 10 channels per trigger         | *&#10004; or &#10008;*
Limit event name length to 200 chars     | *&#10004; or &#10008;*

## Developing the Library

*A section providing information for developers who wish to develop the library*

### Testing

*Any information specific to the library*

### Running tests

*Any additional information specific to the library*

### Deploy to Distribution Mechanism

*Any additional information specific to the library*

## Credits

*It's always nice to give credit to those who inspired the work or who have contributed*

## License

*A statement about the software licence for the library*

**{Example}:**

This code is free to use under the terms of the MIT license.

[channel-names]: https://pusher.com/docs/client_api_guide/client_channels#naming-channels

<!--
##Triggering Events

```go
client := pusher.Client{
    AppId:  app_id,
    Key:    app_key,
    Secret: app_secret,
}

data := map[string]string{"message": "hello world"}

client.Trigger("test_channel", "my_event", data)
```

### Excluding Recipients

In place of the `""` in the last example, we enter the socket_id of the connection we wish to exclude from receiving the event:

```go
client.Trigger("test_channel", "my_event", data, "1234.5678")
```

##Info From All Channels

```go
channelParams := map[string]string{
    "filter_by_prefix": "presence-",
    "info":             "user_count"}

channels, err := client.Channels(channelParams)

fmt.Printf(channels)

// => &{Channels:map[presence-chatroom:{UserCount:4} private-notifications:{UserCount:31}  ]}
```

##Info From One Channel

```go
channelParams := map[string]string{
    "info": "user_count"}

channel, err := client.Channel("presence-chatroom", channelParams)
```

###Gettings Users From Presence Channel

```go
users, err := client.GetChannelUsers("presence-chatroom")
```

## Channel Authentication

### Private Channels

```go
func pusherAuth(res http.ResponseWriter, req *http.Request) {

    params, _ := ioutil.ReadAll(req.Body)
    auth := client.AuthenticateChannel(params)
    fmt.Fprintf(res, auth)

}

func main() {
    http.HandleFunc("/", root)
    http.HandleFunc("/pusher/auth", pusherAuth)

    http.ListenAndServe(":5000", nil)
}

```
### Presence Channels

Like private channels, but one passes in user data to be associated with the member.

```go
params, _ := ioutil.ReadAll(req.Body)

presenceData := pusher.MemberData{
    UserId: "1",
    UserInfo: map[string]string{"twitter": "jamiepatel"}}

auth := client.AuthenticateChannel(params, presenceData)

fmt.Fprintf(res, auth)

```

## Webhooks

```go
body, _ := ioutil.ReadAll(req.Body)

webhook, err := client.Webhook(req.Header, body)

if err != nil {
    fmt.Println("Webhook is invalid :(")
} else {
    fmt.Printf("%+v\n", webhook)
}
```

## Feature Support

Feature                                    | Supported
-------------------------------------------| :-------:
Trigger event on single channel            | *&#10004;*
Trigger event on multiple channels         | *&#10004;*
Limit channels Per trigger to 10           | *&#10004;*
Limit channel/event name length            | *&#10004;*
Validates channel names                    | *&#10004;*
Excluding recipients from events           | *&#10004;*
Authenticating private channels            | *&#10004;*
Authenticating presence channels           | *&#10004;*
Get the list of channels in an application | *&#10004;*
Get the state of a single channel          | *&#10004;*
Get a list of users in a presence channel  | *&#10004;*
WebHook validation                         | *&#10004;*
Can instantiate from URL/ENV               | *&#10004;*
Debugging & Logging                        | *&#10004;*
Cluster configuration                      | *&#10004;*
Timeouts                                   | *&#10004;*
HTTPS                                      | *&#10004;*
HTTP Proxy configuration                   | *&#10008;*
HTTP KeepAlive                             | *&#10008;*


##TODO:

* Finish feature support
* More thorough error-handling.
* Asynchronous requests(?)
* General refactoring

##Running The Tests

    $ go test
 -->
