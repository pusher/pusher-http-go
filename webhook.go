package pusher

type Webhook struct {
	Events      WebhookEvent
	Key, Secret string
}

// func (w *Webhook) IsValid() {

// }
