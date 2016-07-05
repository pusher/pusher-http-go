package pusher

type notificationRequest struct {
	Interests []string `json:"interests"`
	*Notification
}

type Notification struct {
	WebhookURL   string            `json:"webhook_url,omitempty"`
	WebhookLevel string            `json:"webhook_level,omitempty"`
	Apns         *ApnsNotification `json:"apns,omitempty"`
	Gcm          *GcmNotification  `json:"gcm,omitempty"`
}

type ApnsNotification struct {
	Id         string       `json:"apns-id,omitempty"`
	Expiration int64        `json:"expiration,omitempty"`
	Priority   int          `json:"priority,omitempty"`
	CollapseID string       `json:"collapse_id,omitempty"`
	Payload    *ApnsPayload `json:"aps"`
}

type ApnsPayload struct {
	Badge            int        `json:"badge,omitempty"`
	Category         string     `json:"category,omitempty"`
	ContentAvailable int        `json:"content-available,omitempty"`
	URLArgs          []string   `json:"url-args,omitempty"`
	Sound            string     `json:"sound,omitempty"`
	MutableContent   int        `json:"mutable-content,omitempty"`
	Alert            *ApnsAlert `json:"alert,omitempty"`
}

type ApnsAlert struct {
	Action       string   `json:"action,omitempty"`
	ActionLocKey string   `json:"action-loc-key,omitempty"`
	Body         string   `json:"body,omitempty"`
	LaunchImage  string   `json:"launch-image,omitempty"`
	LocArgs      []string `json:"loc-args,omitempty"`
	LocKey       string   `json:"loc-key,omitempty"`
	Title        string   `json:"title,omitempty"`
	Subtitle     string   `json:"subtitle,omitempty"`
	TitleLocArgs []string `json:"title-loc-args,omitempty"`
	TitleLocKey  string   `json:"title-loc-key,omitempty"`
}

type GcmNotification struct {
	CollapseKey           string                 `json:"collapse_key,omitempty"`
	Priority              string                 `json:"priority,omitempty"`
	ContentAvailable      bool                   `json:"content_available,omitempty"`
	DelayWhileIdle        bool                   `json:"delay_while_idle,omitempty"`
	TimeToLive            *uint                  `json:"time_to_live,omitempty"`
	RestrictedPackageName string                 `json:"restricted_package_name,omitempty"`
	DryRun                bool                   `json:"dry_run,omitempty"`
	Data                  map[string]interface{} `json:"data,omitempty"`
	Payload               *GcmPayload            `json:"notification,omitempty"`
}

type GcmPayload struct {
	Title        string `json:"title,omitempty"`
	Body         string `json:"body,omitempty"`
	Icon         string `json:"icon,omitempty"`
	Sound        string `json:"sound,omitempty"`
	Badge        string `json:"badge,omitempty"`
	Tag          string `json:"tag,omitempty"`
	Color        string `json:"color,omitempty"`
	ClickAction  string `json:"click_action,omitempty"`
	BodyLocKey   string `json:"body_loc_key,omitempty"`
	BodyLocArgs  string `json:"body_loc_args,omitempty"`
	TitleLocArgs string `json:"title_loc_args,omitempty"`
	TitleLocKey  string `json:"title_loc_key,omitempty"`
}

type NotificationResponse struct {
	Description         string `json:"description,omitempty"`
	MessageId           string `json:"id,omitempty"`
	SentDeviceToken     string `json:"sent_device_token,omitempty"`
	Success             bool   `json:"success"`
	Platform            string `json:"platform"`
	ReceivedDeviceToken string `json:"received_device_token,omitempty"`
}
