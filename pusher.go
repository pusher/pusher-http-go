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
)

const (
	PUSHER_APP_KEY_HEADER   = "X-Pusher-Key"
	PUSHER_SIGNATURE_HEADER = "X-Pusher-Signature"
)

type Pusher struct {
	appID, key, secret string
	dispatcher
	Options
}

func (p *Pusher) Trigger(channel string, eventName string, data interface{}) (*TriggerResponse, error) {
	return p.trigger(&event{
		Channels: []string{channel},
		Name:     eventName,
		Data:     data,
	})
}

func (p *Pusher) TriggerMulti(channels []string, eventName string, data interface{}) (*TriggerResponse, error) {
	return p.trigger(&event{
		Channels: channels,
		Name:     eventName,
		Data:     data,
	})
}

func (p *Pusher) TriggerExclusive(channel string, eventName string, data interface{}, socketID string) (*TriggerResponse, error) {
	return p.trigger(&event{
		Channels: []string{channel},
		Name:     eventName,
		Data:     data,
		SocketID: &socketID,
	})
}

func (p *Pusher) TriggerMultiExclusive(channels []string, eventName string, data interface{}, socketID string) (*TriggerResponse, error) {
	return p.trigger(&event{
		Channels: channels,
		Name:     eventName,
		Data:     data,
		SocketID: &socketID,
	})
}

func (p *Pusher) trigger(event *event) (*TriggerResponse, error) {
	var eventJSON []byte

	if len(event.Channels) > 10 {
		return nil, errors.New("You cannot trigger on more than 10 channels at once")
	}

	err := validate.Channels(event.Channels)
	if err != nil {
		return nil, err
	}

	err = validate.SocketID(event.SocketID)
	if err != nil {
		return nil, err
	}

	eventJSON, err = json.Marshal(event)
	if err != nil {
		return nil, err
	}

	params := &requests.Params{Body: eventJSON}

	byteResponse, err := p.sendRequest(p.urlConfig(), p.GetHttpClient(), requests.Trigger, params)
	if err != nil {
		return nil, err
	}

	var response *TriggerResponse
	err = json.Unmarshal(byteResponse, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *Pusher) TriggerBatch(batch []Event) (*TriggerResponse, error) {
	batchJSON, err := json.Marshal(&batchRequest{batch})
	if err != nil {
		return nil, err
	}

	params := &requests.Params{
		Body: batchJSON,
	}

	byteResponse, err := p.sendRequest(p.urlConfig(), p.GetHttpClient(), requests.TriggerBatch, params)
	if err != nil {
		return nil, err
	}

	var response *TriggerResponse
	err = json.Unmarshal(byteResponse, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *Pusher) Channels(additionalQueries map[string]string) (*ChannelList, error) {
	params := &requests.Params{
		Queries: additionalQueries,
	}

	byteResponse, err := p.sendRequest(p.urlConfig(), p.GetHttpClient(), requests.Channels, params)
	if err != nil {
		return nil, err
	}

	var response *ChannelList
	err = json.Unmarshal(byteResponse, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *Pusher) Channel(name string, additionalQueries map[string]string) (*Channel, error) {
	params := &requests.Params{
		Channel: name,
		Queries: additionalQueries,
	}

	byteResponse, err := p.sendRequest(p.urlConfig(), p.GetHttpClient(), requests.Channel, params)
	if err != nil {
		return nil, err
	}

	var response *Channel
	err = json.Unmarshal(byteResponse, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *Pusher) ChannelUsers(name string) (*UserList, error) {
	params := &requests.Params{
		Channel: name,
	}

	byteResponse, err := p.sendRequest(p.urlConfig(), p.GetHttpClient(), requests.ChannelUsers, params)
	if err != nil {
		return nil, err
	}

	var response *UserList
	err = json.Unmarshal(byteResponse, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *Pusher) AuthenticatePrivateChannel(body []byte) ([]byte, error) {
	return p.authenticate(&authentications.PrivateChannel{Body: body})
}

func (p *Pusher) AuthenticatePresenceChannel(body []byte, member authentications.Member) ([]byte, error) {
	return p.authenticate(&authentications.PresenceChannel{Body: body, Member: member})
}

func (p *Pusher) authenticate(request authentications.Request) ([]byte, error) {
	unsigned, err := request.StringToSign()
	if err != nil {
		return nil, err
	}

	authSignature := signatures.HMAC(unsigned, p.secret)
	responseMap := map[string]string{
		"auth": fmt.Sprintf("%s:%s", p.key, authSignature),
	}

	userData, err := request.UserData()
	if err != nil {
		return nil, err
	}

	if userData != "" {
		responseMap["channel_data"] = userData
	}

	return json.Marshal(responseMap)
}

func (p *Pusher) Notify(interests []string, notification *Notification) (*NotifyResponse, error) {
	if len(interests) == 0 {
		return nil, errors.New("The interests slice must not be empty")
	}

	req := &notificationRequest{
		Interests:    interests,
		Notification: notification,
	}

	config := &urlConfig{
		appID:  p.appID,
		key:    p.key,
		secret: p.secret,
		host:   p.GetNotificationHost(),
		scheme: p.GetScheme(),
	}

	body, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	params := &requests.Params{
		Body: body,
	}

	byteResponse, err := p.sendRequest(config, p.GetHttpClient(), requests.NativePush, params)
	if err != nil {
		return nil, err
	}

	var response *NotifyResponse
	err = json.Unmarshal(byteResponse, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *Pusher) Webhook(header http.Header, body []byte) (*Webhook, error) {
	for _, token := range header[PUSHER_APP_KEY_HEADER] {
		if token == p.key && signatures.CheckHMAC(header.Get(PUSHER_SIGNATURE_HEADER), p.secret, body) {
			return newWebhook(body)
		}
	}

	return nil, errors.New("Invalid webhook")
}

func (p *Pusher) urlConfig() *urlConfig {
	return &urlConfig{
		appID:  p.appID,
		key:    p.key,
		secret: p.secret,
		host:   p.GetHost(),
		scheme: p.GetScheme(),
	}
}
