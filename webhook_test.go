package pusher

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func setUpClient() Client {
	return Client{AppId: "id", Key: "key", Secret: "secret"}
}

func TestClientWebhookValidation(t *testing.T) {
	client := setUpClient()
	header := make(http.Header)
	header["X-Pusher-Key"] = []string{"key"}
	header["X-Pusher-Signature"] = []string{"2677ad3e7c090b2fa2c0fb13020d66d5420879b8316eb356a2d60fb9073bc778"}
	body := []byte("{\"hello\":\"world\"}")
	webhook, err := client.Webhook(header, body)
	assert.NotNil(t, webhook)
	assert.Nil(t, err)
}

func TestWebhookImproperKeyCase(t *testing.T) {
	client := setUpClient()
	bad_header := make(http.Header)
	bad_header["X-Pusher-Key"] = []string{"narr you're going down!"}
	bad_header["X-Pusher-Signature"] = []string{"2677ad3e7c090b2fa2c0fb13020d66d5420879b8316eb356a2d60fb9073bc778"}
	bad_body := []byte("{\"hello\":\"world\"}")

	bad_webhook, err := client.Webhook(bad_header, bad_body)
	assert.Nil(t, bad_webhook)
	assert.Error(t, err)
}

func TestWebhookImproperSignatureCase(t *testing.T) {
	client := setUpClient()
	bad_header := make(http.Header)
	bad_header["X-Pusher-Key"] = []string{"key"}
	bad_header["X-Pusher-Signature"] = []string{"2677ad3e7c090i'mgonnagetyaeb356a2d60fb9073bc778"}
	bad_body := []byte("{\"hello\":\"world\"}")

	bad_webhook, err := client.Webhook(bad_header, bad_body)
	assert.Nil(t, bad_webhook)
	assert.Error(t, err)
}

func TestWebhookUnmarshalling(t *testing.T) {
	body := []byte("{\"time_ms\":1427233518933,\"events\":[{\"name\":\"client_event\",\"channel\":\"private-channel\",\"event\":\"client-yolo\",\"data\":\"{\\\"yolo\\\":\\\"woot\\\"}\",\"socket_id\":\"44610.7511910\"}]}")
	result, err := unmarshalledWebhook(body)
	expected := &Webhook{
		TimeMs: 1427233518933,
		Events: []WebhookEvent{
			WebhookEvent{
				Name:     "client_event",
				Channel:  "private-channel",
				Event:    "client-yolo",
				Data:     "{\"yolo\":\"woot\"}",
				SocketId: "44610.7511910",
			},
		},
	}

	assert.Equal(t, expected, result)
	assert.NoError(t, err)
}
