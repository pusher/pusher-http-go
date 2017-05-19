package pusher

import "testing"

func TestPushNotificationValidate(t *testing.T) {
	pnNoPayload := PushNotification{
		Interests:    []string{"testInterest"},
		WebhookURL:   "testURL",
		WebhookLevel: WebhookLvlDebug,
	}

	err := pnNoPayload.validate()
	if err == nil {
		t.Error("Invalid PushNotification with no GCM, FCM or APNS payload did not return an error")
	}

	pnNoInterests := PushNotification{
		Interests:    []string{},
		WebhookURL:   "testURL",
		WebhookLevel: WebhookLvlDebug,
		GCM:          []byte(`hello`),
	}

	err = pnNoInterests.validate()
	if err == nil {
		t.Error("Invalid PushNotification with no Interests did not return an error")
	}

	pnManyInterests := PushNotification{
		Interests:    []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
		WebhookURL:   "testURL",
		WebhookLevel: WebhookLvlDebug,
		GCM:          []byte(`hello`),
	}

	err = pnManyInterests.validate()
	if err == nil {
		t.Error("Invalid PushNotification with 10 < Interests did not return an error")
	}

	pnValid := PushNotification{
		Interests:    []string{"testInterest"},
		WebhookURL:   "testURL",
		WebhookLevel: WebhookLvlDebug,
		GCM:          []byte(`hello`),
	}

	err = pnValid.validate()
	if err != nil {
		t.Error(err)
	}
}
