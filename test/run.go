package main

import (
	"fmt"
	"github.com/user/pusher-http-go"
)

func main() {

	client := pusher.Client{
		AppId:  "92870",
		Key:    "235ba307396082ec6719",
		Secret: "eb2da264db8ec0d07065"}

	data := map[string]string{"message": "hello world"}

	err, res := client.Trigger("test_channel", "my_event", data)

	if err != {
		panic(err)
	}

	fmt.Println(res)

}
