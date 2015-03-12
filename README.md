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

fmt.Println(channel)

//{
//  occupied: true,
//  user_count: 42,
//  subscription_count: 42
// }
```
