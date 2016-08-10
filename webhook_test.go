package pusher

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClientWebhookValidation(t *testing.T) {
	client := New("id", "key", "secret")
	header := make(http.Header)
	header["X-Pusher-Key"] = []string{"key"}
	header["X-Pusher-Signature"] = []string{"2677ad3e7c090b2fa2c0fb13020d66d5420879b8316eb356a2d60fb9073bc778"}
	body := []byte("{\"hello\":\"world\"}")
	webhook, err := client.Webhook(header, body)
	assert.NotNil(t, webhook)
	assert.Nil(t, err)
}

func TestWebhookImproperKeyCase(t *testing.T) {
	client := New("id", "key", "secret")
	badHeader := make(http.Header)
	badHeader["X-Pusher-Key"] = []string{"narr you're going down!"}
	badHeader["X-Pusher-Signature"] = []string{"2677ad3e7c090b2fa2c0fb13020d66d5420879b8316eb356a2d60fb9073bc778"}
	badBody := []byte("{\"hello\":\"world\"}")

	badWebhook, err := client.Webhook(badHeader, badBody)
	assert.Nil(t, badWebhook)
	assert.Error(t, err)
}

func TestWebhookImproperSignatureCase(t *testing.T) {
	client := New("id", "key", "secret")
	badHeader := make(http.Header)
	badHeader["X-Pusher-Key"] = []string{"key"}
	badHeader["X-Pusher-Signature"] = []string{"2677ad3e7c090i'mgonnagetyaeb356a2d60fb9073bc778"}
	badBody := []byte("{\"hello\":\"world\"}")

	badWebhook, err := client.Webhook(badHeader, badBody)
	assert.Nil(t, badWebhook)
	assert.Error(t, err)
}

func TestWebhookNoSignature(t *testing.T) {
	client := New("id", "key", "secret")
	badHeader := make(http.Header)
	badHeader["X-Pusher-Key"] = []string{"key"}
	badBody := []byte("{\"hello\":\"world\"}")

	badWebhook, err := client.Webhook(badHeader, badBody)
	assert.Nil(t, badWebhook)
	assert.Error(t, err)
}

func TestWebhookUnmarshalling(t *testing.T) {
	body := []byte("{\"time_ms\":1427233518933,\"events\":[{\"name\":\"client_event\",\"channel\":\"private-channel\",\"event\":\"client-yolo\",\"data\":\"{\\\"yolo\\\":\\\"woot\\\"}\",\"socket_id\":\"44610.7511910\"}]}")
	result, err := newWebhook(body)
	expected := &Webhook{
		TimeMs: 1427233518933,
		Events: []WebhookEvent{{
			Name:     "client_event",
			Channel:  "private-channel",
			Event:    "client-yolo",
			Data:     "{\"yolo\":\"woot\"}",
			SocketId: "44610.7511910",
		}},
	}

	assert.Equal(t, expected, result)
	assert.NoError(t, err)
}
