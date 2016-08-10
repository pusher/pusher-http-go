package errors

import (
	"fmt"
)

const ERROR_TAG = "[pusher-http-go]"

func New(message string) error {
	return fmt.Errorf("%s: %s", ERROR_TAG, message)
}
