package requests

var (
	Trigger = &Request{
		Method:      "POST",
		PathPattern: "/apps/%s/events",
	}

	TriggerBatch = &Request{
		Method:      "POST",
		PathPattern: "/apps/%s/batch_events",
	}

	Channels = &Request{
		Method:      "GET",
		PathPattern: "/apps/%s/channels",
	}

	Channel = &Request{
		Method:      "GET",
		PathPattern: "/apps/%s/channels/%s",
	}

	ChannelUsers = &Request{
		Method:      "GET",
		PathPattern: "/apps/%s/channels/%s/users",
	}

	NativePush = &Request{
		Method:      "POST",
		PathPattern: "/server_api/v1/apps/%s/notifications",
	}
)
