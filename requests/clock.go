package requests

import (
	"strconv"
	"time"
)

type clock interface {
	Now() string
}

type defaultClock struct{}

func (d defaultClock) Now() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

var authClock clock

func init() {
	authClock = defaultClock{}
}
