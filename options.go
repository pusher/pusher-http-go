package pusher

import (
	"fmt"
	"net/http"
	"time"
)

type Options struct {
	Host, Cluster, NotificationHost string // host or host:port pair
	Secure                          bool   // true for HTTPS
	HttpClient                      *http.Client
}

const (
	DEFAULT_HOST              = "api.pusherapp.com"
	DEFAULT_NOTIFICATION_HOST = "nativepush-cluster1.pusher.com"
)

func (o *Options) GetHost() string {
	if o.Host == "" {
		if cluster := o.Cluster; o.Cluster != "" {
			o.Cluster = fmt.Sprintf("api-%s.pusher.com", cluster)
		} else {
			o.Host = DEFAULT_HOST
		}
	}
	return o.Host
}

func (o *Options) GetScheme() string {
	if o.Secure {
		return "https"
	}
	return "http"
}

func (o *Options) GetNotificationHost() string {
	if o.NotificationHost == "" {
		o.NotificationHost = DEFAULT_NOTIFICATION_HOST
	}
	return o.NotificationHost
}

func (o *Options) GetHttpClient() *http.Client {
	if o.HttpClient == nil {
		o.HttpClient = &http.Client{Timeout: time.Second * 5}
	}
	return o.HttpClient
}
