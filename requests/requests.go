package requests

var Trigger = &Request{
	Method:      "POST",
	PathPattern: "/apps/%s/events",
}

var TriggerBatch = &Request{
	Method:      "POST",
	PathPattern: "/apps/%s/batch_events",
}

var Channels = &Request{
	Method:      "GET",
	PathPattern: "/apps/%s/channels",
}

var Channel = &Request{
	Method:      "GET",
	PathPattern: "/apps/%s/channels/%s",
}

var ChannelUsers = &Request{
	Method:      "GET",
	PathPattern: "/apps/%s/channels/%s/users",
}

var NativePush = &Request{
	Method:      "POST",
	PathPattern: "/server_api/v1/apps/%s/notifications",
}
