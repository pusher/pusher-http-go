package pusher

import "testing"

func TestNotificationRequestValidate(t *testing.T) {
	testPayload := PushNotification{
		WebhookURL: "testURL",
	}

	pNReqNoPayload := notificationRequest{
		[]string{"testInterest"},
		&testPayload,
	}

	err := pNReqNoPayload.validate()
	if err == nil {
		t.Error("Invalid PushNotification with no GCM, FCM or APNS payload did not return an error")
	}

	testPayload = PushNotification{
		WebhookURL: "testURL",
		GCM:        []byte(`hello`),
	}

	pnReqNoInterests := notificationRequest{
		[]string{},
		&testPayload,
	}

	err = pnReqNoInterests.validate()
	if err == nil {
		t.Error("Invalid PushNotification with no Interests did not return an error")
	}

	testPayload = PushNotification{
		WebhookURL: "testURL",
		GCM:        []byte(`hello`),
	}

	pnReqManyInterests := notificationRequest{
		[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
		&testPayload,
	}

	err = pnReqManyInterests.validate()
	if err == nil {
		t.Error("Invalid PushNotification with 10 < Interests did not return an error")
	}

	testPayload = PushNotification{
		WebhookURL: "testURL",
		GCM:        []byte(`hello`),
	}

	pnReqValid := notificationRequest{
		[]string{"testInterest"},
		&testPayload,
	}

	err = pnReqValid.validate()
	if err != nil {
		t.Error(err)
	}
}
