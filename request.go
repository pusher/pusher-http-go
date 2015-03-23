package pusher

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func request(method, url string, body []byte) (error, []byte) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	// fmt.Println(resp)

	// fmt.Printf("%+v\n", resp)

	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	// fmt.Println(string(resp_body))

	return process_response(resp.StatusCode, resp_body)
}

func process_response(status int, resp_body []byte) (error, []byte) {
	if status == 200 {
		return nil, resp_body
	} else {
		message := fmt.Sprintf("Status Code: %s - %s", strconv.Itoa(status), string(resp_body))
		err := errors.New(message)
		return err, nil
	}
}
