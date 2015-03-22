package pusher

import "testing"

func TestTriggerRequestUrl(t *testing.T) {

	expected := "http://api.pusherapp.com/apps/3/events?auth_key=278d425bdf160c739803&auth_signature=da454824c97ba181a32ccc17a72625ba02771f50b50e1e7430e47a1f3f457e6c&auth_timestamp=1353088179&auth_version=1.0&body_md5=ec365a775a4cd0599faeb73354201b6f"

	payload := []byte("{\"name\":\"foo\",\"channels\":[\"project-3\"],\"data\":\"{\\\"some\\\":\\\"data\\\"}\"}")

	result := CreateRequestUrl("POST", "/apps/3/events", "278d425bdf160c739803", "7ad3773142a6692b25b8", "1353088179", payload, nil)

	if result != expected {
		t.Error("Expected "+expected+", got", result)
	}

}

func TestGetAllChannels(t *testing.T) {

	expected := "http://api.pusherapp.com/apps/102015/channels?auth_key=d41a439c438a100756f5&auth_signature=4d8a02edcc8a758b0162cd6da690a9a45fb8ae326a276dca1e06a0bc42796c11&auth_timestamp=1427034994&auth_version=1.0&filter_by_prefix=presence-&info=user_count"

	additional_queries := map[string]string{
		"filter_by_prefix": "presence-",
		"info":             "user_count"}

	result := CreateRequestUrl("GET", "/apps/102015/channels", "d41a439c438a100756f5", "4bf35003e819bb138249", "1427034994", nil, additional_queries)

	if result != expected {
		t.Error("Expected "+expected+", got", result)
	}

}
