package pusher

type notificationRequest struct {
	Interests []string `json:"interests"`
	*Notification
}

type Notification struct {
	WebhookURL   string                `json:"webhook_url,omitempty"`
	WebhookLevel string                `json:"webhook_level,omitempty"`
	Apns         *ApnsPushNotification `json:"apns,omitempty"`
	Gcm          *GcmPushNotification  `json:"gcm,omitempty"`
}

type ApnsPushNotification struct {
	Id         string                 `json:"apns-id,omitempty"`
	Expiration int64                  `json:"expiration,omitempty"`
	Priority   int                    `json:"priority,omitempty"`
	CollapseID string                 `json:"collapse_id,omitempty"`
	Aps        *ApnsAps               `json:"aps,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

type ApnsAps struct {
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

type GcmPushNotification struct {
	CollapseKey           string                 `json:"collapse_key,omitempty"`
	Priority              string                 `json:"priority,omitempty"`
	ContentAvailable      bool                   `json:"content_available,omitempty"`
	DelayWhileIdle        bool                   `json:"delay_while_idle,omitempty"`
	TimeToLive            *uint                  `json:"time_to_live,omitempty"`
	RestrictedPackageName string                 `json:"restricted_package_name,omitempty"`
	DryRun                bool                   `json:"dry_run,omitempty"`
	Data                  map[string]interface{} `json:"data,omitempty"`
	Notification          *GcmNotification       `json:"notification,omitempty"`
}

type GcmNotification struct {
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

type NotifyResponse struct {
	NumSubscribers int `json:"number_of_subscribers"`
}
