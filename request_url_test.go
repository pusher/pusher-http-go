package pusher

import "testing"

func logTestResult(expected, result string, t *testing.T) {
	if result != expected {
		t.Error("Expected "+expected+", got", result)
	}
}

func TestTriggerRequestUrl(t *testing.T) {
	expected := "http://api.pusherapp.com/apps/3/events?auth_key=278d425bdf160c739803&auth_signature=da454824c97ba181a32ccc17a72625ba02771f50b50e1e7430e47a1f3f457e6c&auth_timestamp=1353088179&auth_version=1.0&body_md5=ec365a775a4cd0599faeb73354201b6f"
	payload := []byte("{\"name\":\"foo\",\"channels\":[\"project-3\"],\"data\":\"{\\\"some\\\":\\\"data\\\"}\"}")
	result := CreateRequestUrl("POST", "/apps/3/events", "278d425bdf160c739803", "7ad3773142a6692b25b8", "1353088179", payload, nil)
	logTestResult(expected, result, t)
}

func TestGetAllChannelsUrl(t *testing.T) {
	expected := "http://api.pusherapp.com/apps/102015/channels?auth_key=d41a439c438a100756f5&auth_signature=4d8a02edcc8a758b0162cd6da690a9a45fb8ae326a276dca1e06a0bc42796c11&auth_timestamp=1427034994&auth_version=1.0&filter_by_prefix=presence-&info=user_count"
	additional_queries := map[string]string{"filter_by_prefix": "presence-", "info": "user_count"}
	result := CreateRequestUrl("GET", "/apps/102015/channels", "d41a439c438a100756f5", "4bf35003e819bb138249", "1427034994", nil, additional_queries)

	logTestResult(expected, result, t)

}

func TestGetAllChannelsWithOneAdditionalParamUrl(t *testing.T) {
	expected := "http://api.pusherapp.com/apps/102015/channels?auth_key=d41a439c438a100756f5&auth_signature=b540383af4582af5fbb5df7be5472d54bd0838c9c2021c7743062568839e6f97&auth_timestamp=1427036577&auth_version=1.0&filter_by_prefix=presence-"
	additional_queries := map[string]string{"filter_by_prefix": "presence-"}
	result := CreateRequestUrl("GET", "/apps/102015/channels", "d41a439c438a100756f5", "4bf35003e819bb138249", "1427036577", nil, additional_queries)

	logTestResult(expected, result, t)

}

func TestGetAllChannelsWithNoParamsUrl(t *testing.T) {
	expected := "http://api.pusherapp.com/apps/102015/channels?auth_key=d41a439c438a100756f5&auth_signature=df89248f87f6e6d028925e0b04d60f316527a865992ace6936afa91281d8bef0&auth_timestamp=1427036787&auth_version=1.0"
	additional_queries := map[string]string{}
	result := CreateRequestUrl("GET", "/apps/102015/channels", "d41a439c438a100756f5", "4bf35003e819bb138249", "1427036787", nil, additional_queries)

	logTestResult(expected, result, t)

}

func TestGetChannelUrl(t *testing.T) {
	expected := "http://api.pusherapp.com/apps/102015/channels/presence-session-d41a439c438a100756f5-4bf35003e819bb138249-ROpCFmgFhXY?auth_key=d41a439c438a100756f5&auth_signature=f93ceb31f396aef336226efe512aaf339bd5e39c7c2c04b81cc8681dc16ee785&auth_timestamp=1427053326&auth_version=1.0&info=user_count,subscription_count"
	additional_queries := map[string]string{"info": "user_count,subscription_count"}
	result := CreateRequestUrl("GET", "/apps/102015/channels/presence-session-d41a439c438a100756f5-4bf35003e819bb138249-ROpCFmgFhXY", "d41a439c438a100756f5", "4bf35003e819bb138249", "1427053326", nil, additional_queries)

	logTestResult(expected, result, t)

}

func TestGetUsersUrl(t *testing.T) {
	expected := "http://api.pusherapp.com/apps/102015/channels/presence-session-d41a439c438a100756f5-4bf35003e819bb138249-nYJLy67qh52/users?auth_key=d41a439c438a100756f5&auth_signature=207feaf4e8efeb24e5f148011704251bf90e2059a5f97a3eb52d06178b11feca&auth_timestamp=1427053709&auth_version=1.0"
	result := CreateRequestUrl("GET", "/apps/102015/channels/presence-session-d41a439c438a100756f5-4bf35003e819bb138249-nYJLy67qh52/users", "d41a439c438a100756f5", "4bf35003e819bb138249", "1427053709", nil, nil)
	logTestResult(expected, result, t)

}
