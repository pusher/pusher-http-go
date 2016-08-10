package validate

import (
	"github.com/pusher/pusher-http-go/errors"
	"regexp"
)

var socketIDValidationRegex = regexp.MustCompile(`\A\d+\.\d+\z`)

func SocketID(socketID *string) (err error) {
	if (socketID == nil) || socketIDValidationRegex.MatchString(*socketID) {
		return
	}
	return errors.New("socket_id invalid")
}
