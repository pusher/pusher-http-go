package pusher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	AppId, Key, Secret string
}

type Body struct {
	Name     string   `json:"name"`
	Channels []string `json:"channels"`
	Data     string   `json:"data"`
}

func (c *Client) Trigger(channels []string, event string, _data map[string]string) {
	data, _ := json.Marshal(_data)

	payload := c.jsonize(&Body{
		Name:     event,
		Channels: channels,
		Data:     string(data)})

	q := Query{"POST", c.event_path(), c.Key, c.Secret, payload}

	c.post(q.generate(), payload)
}

func (c *Client) jsonize(body *Body) []byte {
	json, _ := json.Marshal(body)
	return json
}

func (c *Client) event_path() string {
	return "/apps/" + c.AppId + "/events"
}

func (c *Client) post(url string, body []byte) {

	fmt.Println(url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	resp_body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(resp_body))
}
