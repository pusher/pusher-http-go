package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func main() {

	app_id := "92870"
	key := "235ba307396082ec6719"
	secret := "eb2da264db8ec0d07065"

	base_url := "http://api.pusherapp.com/apps/" + app_id + "/events"

	channel := "test_channel"
	event := "my_event"

	_data := map[string]string{"message": "hello world"}

	data, _ := json.Marshal(_data)

	type Body struct {
		Name     string   `json:"name"`
		Channels []string `json:"channels"`
		Data     string   `json:"data"`
	}

	_body := &Body{
		Name:     event,
		Channels: []string{channel},
		Data:     string(data)}

	body, _ := json.Marshal(_body)

	_body_md5 := md5.New()
	_body_md5.Write([]byte(body))

	body_md5 := hex.EncodeToString(_body_md5.Sum(nil))

	auth_timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	fmt.Println(auth_timestamp)

	auth_version := "1.0"

	method := "POST\n"

	url_path := "/apps/" +
		app_id +
		"/events\n"

	query_string := "auth_key=" + key + "&" +
		"auth_timestamp=" + auth_timestamp + "&" +
		"auth_version=" + auth_version + "&" +
		"body_md5=" + body_md5

	to_sign := method + url_path + query_string

	_auth_signature := hmac.New(sha256.New, []byte(secret))
	_auth_signature.Write([]byte(to_sign))
	auth_signature := hex.EncodeToString(_auth_signature.Sum(nil))

	url := base_url + "?" + query_string + "&auth_signature=" + auth_signature

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
