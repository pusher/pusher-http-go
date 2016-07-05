package pusher

import (
	"encoding/json"
	"fmt"
	"github.com/pusher/pusher-http-go/authentications"
	"github.com/pusher/pusher-http-go/errors"
	"github.com/pusher/pusher-http-go/requests"
	"github.com/pusher/pusher-http-go/signatures"
	"github.com/pusher/pusher-http-go/validate"
	"net/http"
	"net/url"
	"time"
)

const (
	NATIVE_NOTIFICATIONS_HOST = "yolo.ngrok.io"
)

type Pusher struct {
	appID, key, secret string
	dispatcher
	Options
}

type Options struct {
	Host       string // host or host:port pair
	Secure     bool   // true for HTTPS
	Cluster    string
	HttpClient *http.Client
}

func (p *Pusher) Trigger(channel string, eventName string, data interface{}) (*TriggerResponse, error) {
	event := &event{
		Channels: []string{channel},
		Name:     eventName,
		Data:     data,
	}
	return p.trigger(event)
}

func (p *Pusher) TriggerMulti(channels []string, eventName string, data interface{}) (*TriggerResponse, error) {
	event := &event{
		Channels: channels,
		Name:     eventName,
		Data:     data,
	}
	return p.trigger(event)
}

func (p *Pusher) TriggerExclusive(channel string, eventName string, data interface{}, socketID string) (*TriggerResponse, error) {
	event := &event{
		Channels: []string{channel},
		Name:     eventName,
		Data:     data,
		SocketID: &socketID,
	}
	return p.trigger(event)
}

func (p *Pusher) TriggerMultiExclusive(channels []string, eventName string, data interface{}, socketID string) (*TriggerResponse, error) {
	event := &event{
		Channels: channels,
		Name:     eventName,
		Data:     data,
		SocketID: &socketID,
	}
	return p.trigger(event)
}

func (p *Pusher) trigger(event *event) (response *TriggerResponse, err error) {
	var (
		eventJSON    []byte
		byteResponse []byte
	)

	if len(event.Channels) > 10 {
		err = errors.New("You cannot trigger on more than 10 channels at once")
		return
	}

	if err = validate.Channels(event.Channels); err != nil {
		return
	}

	if err = validate.SocketID(event.SocketID); err != nil {
		return
	}

	if eventJSON, err = json.Marshal(event); err != nil {
		return
	}

	params := &requests.Params{
		Body: eventJSON,
	}

	if byteResponse, err = p.sendRequest(p.urlConfig(), p.httpClient(), requests.Trigger, params); err != nil {
		return
	}

	err = json.Unmarshal(byteResponse, &response)
	return
}

func (p *Pusher) TriggerBatch(batch []Event) (response *TriggerResponse, err error) {
	var (
		byteResponse []byte
		batchJSON    []byte
	)

	if batchJSON, err = json.Marshal(&batchRequest{batch}); err != nil {
		return
	}

	params := &requests.Params{
		Body: batchJSON,
	}

	if byteResponse, err = p.sendRequest(p.urlConfig(), p.httpClient(), requests.TriggerBatch, params); err != nil {
		return
	}

	err = json.Unmarshal(byteResponse, &response)
	return
}

func (p *Pusher) Channels(additionalQueries map[string]string) (response *ChannelList, err error) {
	var byteResponse []byte

	params := &requests.Params{
		Queries: additionalQueries,
	}

	if byteResponse, err = p.sendRequest(p.urlConfig(), p.httpClient(), requests.Channels, params); err != nil {
		return
	}

	err = json.Unmarshal(byteResponse, &response)
	return
}

func (p *Pusher) Channel(name string, additionalQueries map[string]string) (response *Channel, err error) {
	var byteResponse []byte

	params := &requests.Params{
		Channel: name,
		Queries: additionalQueries,
	}

	if byteResponse, err = p.sendRequest(p.urlConfig(), p.httpClient(), requests.Channel, params); err != nil {
		return
	}

	err = json.Unmarshal(byteResponse, &response)
	return
}

func (p *Pusher) ChannelUsers(name string) (response *UserList, err error) {
	var byteResponse []byte

	params := &requests.Params{
		Channel: name,
	}

	if byteResponse, err = p.sendRequest(p.urlConfig(), p.httpClient(), requests.ChannelUsers, params); err != nil {
		return
	}

	err = json.Unmarshal(byteResponse, &response)
	return
}

func (p *Pusher) AuthenticatePrivateChannel(body []byte) (response []byte, err error) {
	return p.authenticate(&authentications.PrivateChannel{Body: body})
}

func (p *Pusher) AuthenticatePresenceChannel(body []byte, member authentications.Member) (response []byte, err error) {
	return p.authenticate(&authentications.PresenceChannel{Body: body, Member: member})
}

func (p *Pusher) authenticate(request authentications.Request) (response []byte, err error) {
	var unsigned string
	if unsigned, err = request.StringToSign(); err != nil {
		return
	}
	authSignature := signatures.HMAC(unsigned, p.secret)

	responseMap := map[string]string{
		"auth": fmt.Sprintf("%s:%s", p.key, authSignature),
	}
	var userData string
	if userData, err = request.UserData(); err != nil {
		return
	}
	if userData != "" {
		responseMap["channel_data"] = userData
	}
	return json.Marshal(responseMap)
}

func (p *Pusher) Notify(interests []string, notification *Notification) (response []NotificationResponse, err error) {
	if len(interests) != 1 {
		err = errors.New("Currently only sending to one interest is supported.")
		return
	}

	if notification.Gcm != nil {
		ttl := notification.Gcm.TimeToLive
		if ttl != nil && (*ttl < uint(0) || uint(2419200) < *ttl) {
			err = errors.New("GCM's TimeToLive field must be an integer between 0 and 2419200 (4 weeks)")
			return
		}

		if notification.Gcm.Payload != nil {
			payload := notification.Gcm.Payload
			if len(payload.Title) == 0 {
				err = errors.New("Notification title is a required field for GCM")
				return
			}

			if len(payload.Icon) == 0 {
				err = errors.New("Notification icon is a required field for GCM")
				return
			}

		}

		if notification.WebhookURL != "" {
			if _, err = url.Parse(notification.WebhookURL); err != nil {
				return
			}
		}
	}

	if !map[string]bool{"INFO": true, "DEBUG": true, "": true}[notification.WebhookLevel] {
		err = errors.New("Webhook level must be either INFO or DEBUG. Blank will default to INFO")
		return
	}

	var body, byteResponse []byte

	req := &notificationRequest{
		Interests:    interests,
		Notification: notification,
	}

	config := &urlConfig{
		key:    p.key,
		secret: p.secret,
		host:   NATIVE_NOTIFICATIONS_HOST,
		appID:  p.appID,
	}

	if body, err = json.Marshal(&req); err != nil {
		return
	}

	fmt.Println(string(body))

	params := &requests.Params{
		Body: body,
	}

	if byteResponse, err = p.sendRequest(config, p.httpClient(), requests.NativePush, params); err != nil {
		return
	}

	err = json.Unmarshal(byteResponse, &response)
	return
}

func (p *Pusher) Webhook(header http.Header, body []byte) (webhook *Webhook, err error) {
	for _, token := range header["X-Pusher-Key"] {
		if token == p.key && signatures.CheckHMAC(header.Get("X-Pusher-Signature"), p.secret, body) {
			return newWebhook(body)
		}
	}
	err = errors.New("Invalid webhook")
	return
}

func (p *Pusher) urlConfig() *urlConfig {
	return &urlConfig{
		key:     p.key,
		secret:  p.secret,
		host:    p.Host,
		cluster: p.Cluster,
		appID:   p.appID,
		secure:  p.Secure,
	}
}

func (p *Pusher) httpClient() *http.Client {
	if p.HttpClient == nil {
		p.HttpClient = &http.Client{Timeout: time.Second * 5}
	}

	return p.HttpClient
}
