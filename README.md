# Pusher Go Library

##Triggering Events

```go
client := pusher.Client{
    AppId:  app_id,
    Key:    app_key,
    Secret: app_secret}

data := map[string]string{"message": "hello world"}

client.Trigger([]string{"test_channel"}, "my_event", data)
```

##Info From All Channels

```go
channelParams := map[string]string{
    "filter_by_prefix": "presence-",
    "info":             "user_count"}

err, channels := client.Channels(channelParams)

fmt.Println(channels)

// => { "channels":
//        { "presence-chatroom":
//            { "user_count": 1 }
//        }
//    }

```

##Info From One Channel

```go
channelParams := map[string]string{
    "info": "user_count"}

err, channel := client.Channel("presence-chatroom", channelParams)
```

###Gettings Users From Presence Channel

```go
err, users := client.GetChannelUsers("presence-chatroom")
```

## Channel Authentication

### Private Channels

```go
func pusher_auth(res http.ResponseWriter, req *http.Request) {

    params, _ := ioutil.ReadAll(req.Body)
    auth := client.AuthenticateChannel(params)
    fmt.Fprintf(res, auth)

}

func main() {
    http.HandleFunc("/", root)
    http.HandleFunc("/pusher/auth", pusher_auth)

    http.ListenAndServe(":5000", nil)
}

```
### Presence Channels

Like private channels, but one passes in user data to be associated with the member.

```go
params, _ := ioutil.ReadAll(req.Body)

presence_data := map[string]string{"user_id": "1"}
auth := client.AuthenticateChannel(params, presence_data)

fmt.Fprintf(res, auth)

```



