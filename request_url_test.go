package pusher

import "testing"

func TestRequestUrl(t *testing.T) {

	var result string

	expected := "http://api.pusherapp.com/apps/3/events?auth_key=278d425bdf160c739803&auth_signature=da454824c97ba181a32ccc17a72625ba02771f50b50e1e7430e47a1f3f457e6c&auth_timestamp=1353088179&auth_version=1.0&body_md5=ec365a775a4cd0599faeb73354201b6f"

	payload := []byte("{\"name\":\"foo\",\"channels\":[\"project-3\"],\"data\":\"{\\\"some\\\":\\\"data\\\"}\"}")

	result = CreateRequestUrl("POST", "/apps/3/events", "278d425bdf160c739803", "7ad3773142a6692b25b8", "1353088179", payload, nil)

	if result != expected {
		t.Error("Expected "+expected+", got", result)
	}

}
