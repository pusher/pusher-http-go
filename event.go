package pusher

type Event struct {
	Name     string   `json:"name"`
	Channels []string `json:"channels"`
	Data     string   `json:"data"`
	SocketId string   `json:"socket_id"`
}

type WebhookEvent struct {
	Name     string `json:"name"`
	Channel  string `json:"channel"`
	Event    string `json:"event"`
	Data     string `json:"data"`
	SocketId string `json:"socket_id"`
}
