package pusher

import (
	"github.com/pusher/pusher-http-go/errors"
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

func NewFromURL(rawURL string) (Client, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	host := u.Host

	matches := pusherPathRegex.FindStringSubmatch(u.Path)
	if len(matches) == 0 {
		return nil, errors.New("No app ID found")
	}

	appID := matches[1]

	if u.User == nil {
		return nil, errors.New("Missing <key>:<secret>")
	}

	key := u.User.Username()

	secret, secretGiven := u.User.Password()
	if !secretGiven {
		return nil, errors.New("Missing <secret>")
	}

	secure := u.Scheme == "https"

	client := NewWithOptions(appID, key, secret, Options{
		Secure: secure,
		Host:   host,
	})

	return client, nil
}

func NewFromEnv(key string) (Client, error) {
	return NewFromURL(os.Getenv(key))
}
