package pusher

import (
	"github.com/pusher/pusher/errors"
	"net/url"
	"os"
	"regexp"
)

func New(appID, key, secret string) Client {
	return &Pusher{
		appID:  appID,
		key:    key,
		secret: secret,
		Options: Options{
			Host:   "api.pusherapp.com",
			Secure: true,
		},
		dispatcher: defaultDispatcher{},
	}
}

func NewWithOptions(appID, key, secret string, options Options) Client {
	return &Pusher{
		appID:      appID,
		key:        key,
		secret:     secret,
		Options:    options,
		dispatcher: defaultDispatcher{},
	}
}

var pusherPathRegex = regexp.MustCompile("^/apps/([0-9]+)$")

func NewFromURL(rawURL string) (client Client, err error) {
	var (
		u           *url.URL
		host        string
		appID       string
		key         string
		secret      string
		secretGiven bool
		secure      bool
	)

	if u, err = url.Parse(rawURL); err != nil {
		return
	}

	host = u.Host

	matches := pusherPathRegex.FindStringSubmatch(u.Path)
	if len(matches) == 0 {
		err = errors.New("No app ID found")
		return
	}
	appID = matches[1]

	if u.User == nil {
		err = errors.New("Missing <key>:<secret>")
		return
	}
	key = u.User.Username()

	if secret, secretGiven = u.User.Password(); !secretGiven {
		err = errors.New("Missing <secret>")
		return
	}

	secure = u.Scheme == "https"

	client = NewWithOptions(appID, key, secret, Options{
		Secure: secure,
		Host:   host,
	})
	return
}

func NewFromEnv(key string) (Client, error) {
	return NewFromURL(os.Getenv(key))
}
