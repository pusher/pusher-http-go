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