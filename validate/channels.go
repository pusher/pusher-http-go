package validate

import (
	"fmt"
	"github.com/pusher/pusher-http-go/errors"
	"regexp"
	s "strings"
)

var channelValidationRegex = regexp.MustCompile("^[-a-zA-Z0-9_=@,.;]+$")

func Channels(channels []string) error {
	var channelErrors []string

	for _, channel := range channels {
		if len(channel) > 200 {
			channelErrors = append(channelErrors, channelTooLong(channel))
			continue
		}

		if !channelValidationRegex.MatchString(channel) {
			channelErrors = append(channelErrors, channelHasIllegalCharacters(channel))
			continue
		}
	}

	if len(channelErrors) > 0 {
		message := s.Join(channelErrors, ". ")
		return errors.New(message)
	}

	return nil
}

func channelTooLong(channel string) string {
	return fmt.Sprintf("%s is over 200 characters", channel)
}

func channelHasIllegalCharacters(channel string) string {
	return fmt.Sprintf("%s has illegal characters", channel)
}
