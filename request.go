package pusher

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Request(method, url string, body []byte) (error, []byte) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(resp_body))

	return nil, resp_body
}
