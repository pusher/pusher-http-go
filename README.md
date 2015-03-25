# Pusher Go Library

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

