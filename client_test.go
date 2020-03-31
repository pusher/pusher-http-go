package pusher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTriggerSuccessCase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		fmt.Fprintf(res, "{}")
		assert.Equal(t, "POST", req.Method)

		expectedBody := "{\"name\":\"test\",\"channels\":[\"test_channel\"],\"data\":\"yolo\"}"
		actualBody, err := ioutil.ReadAll(req.Body)
		assert.Equal(t, expectedBody, string(actualBody))

		assert.Equal(t, "application/json", req.Header["Content-Type"][0])
		assert.NoError(t, err)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := Client{AppID: "id", Key: "key", Secret: "secret", Host: u.Host}
	err := client.Trigger("test_channel", "test", "yolo")
	assert.NoError(t, err)
}

func TestTriggerMultiSuccessCase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		fmt.Fprintf(res, "{}")
		assert.Equal(t, "POST", req.Method)

		expectedBody := "{\"name\":\"test\",\"channels\":[\"test_channel\",\"other_channel\"],\"data\":\"yolo\"}"
		actualBody, err := ioutil.ReadAll(req.Body)
		assert.Equal(t, expectedBody, string(actualBody))

		assert.Equal(t, "application/json", req.Header["Content-Type"][0])
		assert.NoError(t, err)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := Client{AppID: "id", Key: "key", Secret: "secret", Host: u.Host}
	err := client.TriggerMulti([]string{"test_channel", "other_channel"}, "test", "yolo")
	assert.NoError(t, err)
}

func TestTriggerMultiEncryptedRejected(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		t.Fatal("No request should reach the API")
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := Client{
		AppID:                     "id",
		Key:                       "key",
		Secret:                    "secret",
		Host:                      u.Host,
		EncryptionMasterKeyBase64: "ZUhQVldIZzduRkdZVkJzS2pPRkRYV1JyaWJJUjJiMGI=",
	}
	err := client.TriggerMulti([]string{"test_channel", "private-encrypted-other_channel"}, "test", "yolo")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "multiple channels")
	assert.Contains(t, err.Error(), "encrypted channels")
}

func TestGetChannelsSuccessCase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		testJSON := "{\"channels\":{\"presence-session-d41a439c438a100756f5-4bf35003e819bb138249-5cbTiUiPNGI\":{\"user_count\":1},\"presence-session-d41a439c438a100756f5-4bf35003e819bb138249-PbZ5E1pP8uF\":{\"user_count\":1},\"presence-session-d41a439c438a100756f5-4bf35003e819bb138249-oz6iqpSxMwG\":{\"user_count\":1}}}"

		fmt.Fprintf(res, testJSON)
		assert.Equal(t, "GET", req.Method)

	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := Client{AppID: "id", Key: "key", Secret: "secret", Host: u.Host}
	channels, err := client.Channels(nil)
	assert.NoError(t, err)

	expected := &ChannelsList{
		Channels: map[string]ChannelListItem{
			"presence-session-d41a439c438a100756f5-4bf35003e819bb138249-5cbTiUiPNGI": ChannelListItem{UserCount: 1},
			"presence-session-d41a439c438a100756f5-4bf35003e819bb138249-PbZ5E1pP8uF": ChannelListItem{UserCount: 1},
			"presence-session-d41a439c438a100756f5-4bf35003e819bb138249-oz6iqpSxMwG": ChannelListItem{UserCount: 1},
		},
	}
	assert.Equal(t, channels, expected)
}

func TestGetChannelSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		testJSON := "{\"user_count\":1,\"occupied\":true,\"subscription_count\":1}"
		fmt.Fprintf(res, testJSON)

		assert.Equal(t, "GET", req.Method)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := Client{AppID: "id", Key: "key", Secret: "secret", Host: u.Host}
	channel, err := client.Channel("test_channel", nil)
	assert.NoError(t, err)

	expected := &Channel{
		Name:              "test_channel",
		Occupied:          true,
		UserCount:         1,
		SubscriptionCount: 1,
	}
	assert.Equal(t, channel, expected)
}

func TestGetChannelUserSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		testJSON := "{\"users\":[{\"id\":\"red\"},{\"id\":\"blue\"}]}"
		fmt.Fprintf(res, testJSON)

		assert.Equal(t, "GET", req.Method)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := Client{AppID: "id", Key: "key", Secret: "secret", Host: u.Host}
	users, err := client.GetChannelUsers("test_channel")
	assert.NoError(t, err)

	expected := &Users{
		List: []User{User{ID: "red"}, User{ID: "blue"}},
	}
	assert.Equal(t, users, expected)
}

func TestTriggerWithSocketID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		expectedBody := "{\"name\":\"test\",\"channels\":[\"test_channel\"],\"data\":\"yolo\",\"socket_id\":\"1234.12\"}"
		actualBody, err := ioutil.ReadAll(req.Body)
		assert.Equal(t, expectedBody, string(actualBody))
		assert.NoError(t, err)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := Client{AppID: "id", Key: "key", Secret: "secret", Host: u.Host}
	client.TriggerExclusive("test_channel", "test", "yolo", "1234.12")
}

func TestTriggerSocketIDValidation(t *testing.T) {
	client := Client{AppID: "id", Key: "key", Secret: "secret"}
	err := client.TriggerExclusive("test_channel", "test", "yolo", "1234.12:lalala")
	assert.Error(t, err)
}

func TestTriggerBatchSuccess(t *testing.T) {
	expectedBody := `{"batch":[{"channel":"test_channel","name":"test","data":"yolo1"},{"channel":"test_channel","name":"test","data":"yolo2"}]}`
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		fmt.Fprintf(res, "{}")
		assert.Equal(t, "POST", req.Method)

		actualBody, err := ioutil.ReadAll(req.Body)
		assert.Equal(t, expectedBody, string(actualBody))
		assert.Equal(t, "application/json", req.Header["Content-Type"][0])
		assert.Equal(t, "/apps/appid/batch_events", req.URL.Path)
		assert.NoError(t, err)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := Client{AppID: "appid", Key: "key", Secret: "secret", Host: u.Host}
	err := client.TriggerBatch([]Event{
		{"test_channel", "test", "yolo1", nil},
		{"test_channel", "test", "yolo2", nil},
	})

	assert.NoError(t, err)
}

func TestTriggerBatchWithEncryptionMasterKeyNoEncryptedChanSuccess(t *testing.T) {
	expectedBody := `{"batch":[{"channel":"test_channel","name":"test","data":"yolo1"},{"channel":"test_channel","name":"test","data":"yolo2"}]}`
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		fmt.Fprintf(res, "{}")
		assert.Equal(t, "POST", req.Method)

		actualBody, err := ioutil.ReadAll(req.Body)
		assert.Equal(t, expectedBody, string(actualBody))
		assert.Equal(t, "application/json", req.Header["Content-Type"][0])
		assert.Equal(t, "/apps/appid/batch_events", req.URL.Path)
		assert.NoError(t, err)
	}))
	defer server.Close()
	u, _ := url.Parse(server.URL)
	client := Client{AppID: "appid", Key: "key", Secret: "secret", EncryptionMasterKeyBase64: "ZUhQVldIZzduRkdZVkJzS2pPRkRYV1JyaWJJUjJiMGI=", Host: u.Host}
	err := client.TriggerBatch([]Event{
		{"test_channel", "test", "yolo1", nil},
		{"test_channel", "test", "yolo2", nil},
	})

	assert.NoError(t, err)
}

func TestTriggerBatchNoEncryptionMasterKeyWithEncryptedChanFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		t.Fatal("No request should have reached the API")
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := Client{AppID: "appid", Key: "key", Secret: "secret", Host: u.Host}
	err := client.TriggerBatch([]Event{
		{"test_channel", "test", "yolo1", nil},
		{"private-encrypted-test_channel", "test", "yolo2", nil},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "master encryption key")
}

func TestTriggerWithEncryptedChanSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		fmt.Fprintf(res, "{}")
		assert.Equal(t, "POST", req.Method)

		actualBody, err := ioutil.ReadAll(req.Body)
		assert.Contains(t, string(actualBody), "ciphertext")
		assert.Contains(t, string(actualBody), "nonce")
		assert.Equal(t, "application/json", req.Header["Content-Type"][0])
		assert.Equal(t, "/apps/appid/events", req.URL.Path)
		assert.NoError(t, err)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := Client{AppID: "appid", Key: "key", Secret: "secret", EncryptionMasterKeyBase64: "ZUhQVldIZzduRkdZVkJzS2pPRkRYV1JyaWJJUjJiMGI=", Host: u.Host}
	err := client.Trigger("private-encrypted-test_channel", "test", "yolo1")
	assert.NoError(t, err)
}

func TestTriggerBatchWithEncryptedChanSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		fmt.Fprintf(res, "{}")
		assert.Equal(t, "POST", req.Method)

		_, err := ioutil.ReadAll(req.Body)
		assert.Equal(t, "application/json", req.Header["Content-Type"][0])
		assert.Equal(t, "/apps/appid/batch_events", req.URL.Path)
		assert.NoError(t, err)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := Client{AppID: "appid", Key: "key", Secret: "secret", EncryptionMasterKeyBase64: "ZUhQVldIZzduRkdZVkJzS2pPRkRYV1JyaWJJUjJiMGI=", Host: u.Host}
	err := client.TriggerBatch([]Event{
		{"test_channel", "test", "yolo1", nil},
		{"private-encrypted-test_channel", "test", "yolo2", nil},
	})
	assert.NoError(t, err)
}

func TestTriggerInvalidMasterKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		t.Fatal("No HTTP request should have been made")
	}))
	defer server.Close()
	u, _ := url.Parse(server.URL)

	// too short (deprecated)
	client := Client{
		AppID:               "appid",
		Key:                 "key",
		Secret:              "secret",
		Host:                u.Host,
		EncryptionMasterKey: "this is 31 bytes 12345678901234",
	}
	err := client.Trigger("private-encrypted-test_channel", "test", "yolo1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "32 bytes")

	// too long (deprecated)
	client = Client{
		AppID:               "appid",
		Key:                 "key",
		Secret:              "secret",
		Host:                u.Host,
		EncryptionMasterKey: "this is 33 bytes 1234567890123456",
	}
	err = client.Trigger("private-encrypted-test_channel", "test", "yolo1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "32 bytes")

	// both provided
	client = Client{
		AppID:                     "appid",
		Key:                       "key",
		Secret:                    "secret",
		Host:                      u.Host,
		EncryptionMasterKey:       "this is 32 bytes 123456789012345",
		EncryptionMasterKeyBase64: "dGhpcyBpcyAzMiBieXRlcyAxMjM0NTY3ODkwMTIzNDU=",
	}
	err = client.Trigger("private-encrypted-test_channel", "test", "yolo1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "both")

	// too short
	client = Client{
		AppID:                     "appid",
		Key:                       "key",
		Secret:                    "secret",
		Host:                      u.Host,
		EncryptionMasterKeyBase64: "dGhpcyBpcyAzMSBieXRlcyAxMjM0NTY3ODkwMTIzNA==",
	}
	err = client.Trigger("private-encrypted-test_channel", "test", "yolo1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "32 bytes")

	// too long
	client = Client{
		AppID:                     "appid",
		Key:                       "key",
		Secret:                    "secret",
		Host:                      u.Host,
		EncryptionMasterKeyBase64: "dGhpcyBpcyAzMiBieXRlcyAxMjM0NTY3ODkwMTIzNDU2",
	}
	err = client.Trigger("private-encrypted-test_channel", "test", "yolo1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "32 bytes")

	// invalid base64
	client = Client{
		AppID:                     "appid",
		Key:                       "key",
		Secret:                    "secret",
		Host:                      u.Host,
		EncryptionMasterKeyBase64: "dGhp!yBpcyAzMiBieXRlcy#xMjM0NTY3ODkwMTIzNDU=",
	}
	err = client.Trigger("private-encrypted-test_channel", "test", "yolo1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "valid base64")
}

func TestErrorResponseHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(400)
		fmt.Fprintf(res, "Cannot retrieve the user count unless the channel is a presence channel")

	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := Client{AppID: "id", Key: "key", Secret: "secret", Host: u.Host}
	channelParams := map[string]string{"info": "user_count,subscription_count"}
	channel, err := client.Channel("this_is_not_a_presence_channel", channelParams)

	assert.Error(t, err)
	assert.EqualError(t, err, "Status Code: 400 - Cannot retrieve the user count unless the channel is a presence channel")
	assert.Nil(t, channel)
}

func TestRequestTimeouts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		time.Sleep(time.Second * 1)
		// res.WriteHeader(200)
		fmt.Fprintf(res, "{}")
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := Client{AppID: "id", Key: "key", Secret: "secret", Host: u.Host, HTTPClient: &http.Client{Timeout: time.Millisecond * 100}}
	err := client.Trigger("test_channel", "test", "yolo")

	assert.Error(t, err)
}

func TestChannelLengthValidation(t *testing.T) {
	channels := []string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11",
		"12", "13", "14", "15", "16", "17", "18", "19", "20", "21",
		"22", "23", "24", "25", "26", "27", "28", "29", "30", "31",
		"32", "33", "34", "35", "36", "37", "38", "39", "40", "41",
		"42", "43", "44", "45", "46", "47", "48", "49", "50", "51",
		"52", "53", "54", "55", "56", "57", "58", "59", "60", "61",
		"62", "63", "64", "65", "66", "67", "68", "69", "70", "71",
		"72", "73", "74", "75", "76", "77", "78", "79", "80", "81",
		"82", "83", "84", "85", "86", "87", "88", "89", "90", "91",
		"92", "93", "94", "95", "96", "97", "98", "99", "100", "101",
	}

	client := Client{AppID: "id", Key: "key", Secret: "secret"}
	err := client.TriggerMulti(channels, "yolo", "woot")

	assert.EqualError(t, err, "You cannot trigger on more than 100 channels at once")
}

func TestChannelFormatValidation(t *testing.T) {
	channel1 := "w000^$$£@@@"
	var channel2 string
	for i := 0; i <= 202; i++ {
		channel2 += "a"
	}
	client := Client{AppID: "id", Key: "key", Secret: "secret"}
	err1 := client.Trigger(channel1, "yolo", "w00t")

	err2 := client.Trigger(channel2, "yolo", "not 19 forever")

	assert.EqualError(t, err1, "At least one of your channels' names are invalid")

	assert.EqualError(t, err2, "At least one of your channels' names are invalid")

}

func TestDataSizeValidation(t *testing.T) {
	client := Client{AppID: "id", Key: "key", Secret: "secret"}
	var data string
	for i := 0; i <= 10242; i++ {
		data += "a"
	}
	err := client.Trigger("channel", "event", data)

	assert.EqualError(t, err, "Data must be smaller than 10kb")

}

func TestInitialisationFromURL(t *testing.T) {
	url := "http://feaf18a411d3cb9216ee:fec81108d90e1898e17a@api.pusherapp.com/apps/104060"
	client, _ := ClientFromURL(url)
	expectedClient := &Client{Key: "feaf18a411d3cb9216ee", Secret: "fec81108d90e1898e17a", AppID: "104060", Host: "api.pusherapp.com"}
	assert.Equal(t, expectedClient, client)
}

func TestURLInitErrorNoSecret(t *testing.T) {
	url := "http://fec81108d90e1898e17a@api.pusherapp.com/apps"
	client, err := ClientFromURL(url)
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestURLInitHTTPS(t *testing.T) {
	url := "https://key:secret@api.pusherapp.com/apps/104060"
	client, _ := ClientFromURL(url)
	assert.True(t, client.Secure)
}

func TestURLInitErrorNoID(t *testing.T) {
	url := "http://fec81108d90e1898e17a@api.pusherapp.com/apps"
	client, err := ClientFromURL(url)
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestInitialisationFromENV(t *testing.T) {
	os.Setenv("PUSHER_URL", "http://feaf18a411d3cb9216ee:fec81108d90e1898e17a@api.pusherapp.com/apps/104060")
	client, _ := ClientFromEnv("PUSHER_URL")
	expectedClient := &Client{Key: "feaf18a411d3cb9216ee", Secret: "fec81108d90e1898e17a", AppID: "104060", Host: "api.pusherapp.com"}
	assert.Equal(t, expectedClient, client)
}
