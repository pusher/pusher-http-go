package requests

import (
	"net/http"
)

var Trigger = &Request{
	Method:      http.MethodPost,
	PathPattern: "/apps/%s/events",
}

var TriggerBatch = &Request{
	Method:      http.MethodPost,
	PathPattern: "/apps/%s/batch_events",
}

var Channels = &Request{
	Method:      http.MethodGet,
	PathPattern: "/apps/%s/channels",
}

var Channel = &Request{
	Method:      http.MethodGet,
	PathPattern: "/apps/%s/channels/%s",
}

var ChannelUsers = &Request{
	Method:      http.MethodGet,
	PathPattern: "/apps/%s/channels/%s/users",
}

var NativePush = &Request{
	Method:      http.MethodPost,
	PathPattern: "/server_api/v1/apps/%s/notifications",
}
