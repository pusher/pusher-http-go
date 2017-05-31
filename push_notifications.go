package pusher

import "errors"

const (
	WebhookLvlInfo             = "INFO"
	WebhookLvlDebug            = "DEBUG"
	PushNotiHostDefault        = "nativepush-cluster1.pusher.com"
	PushNotifAPIPrefixDefault  = "server_api"
	PushNotifAPIVersionDefault = "v1"
)

// PushNotification is a type for requesting push notifications
type PushNotification struct {
	WebhookURL   string      `json:"webhook_url,omitempty"`
	WebhookLevel string      `json:"webhook_level,omitempty"`
	APNS         interface{} `json:"apns,omitempty"`
	GCM          interface{} `json:"gcm,omitempty"`
	FCM          interface{} `json:"fcm,omitempty"`
}

type notificationRequest struct {
	Interests []string `json:"interests"`
	*PushNotification
}

// validate checks the notificationRequest has 0<Interests<11 and has a
// APNS, GCM or FCM payload
func (pN *notificationRequest) validate() error {
	if 0 == len(pN.Interests) || len(pN.Interests) > 10 {
		return errors.New("Interests must contain between 1 and 10 interests")
	}

	if pN.APNS == nil && pN.GCM == nil && pN.FCM == nil {
		return errors.New("PushNotification must contain a GCM, FCM or APNS payload")
	}

	return nil
}

// NotifyResponse is returned from a successful PushNotification and contain the number of
// subscribers to those interests
type NotifyResponse struct {
	NumSubscribers int `json:"number_of_subscribers"`
}
