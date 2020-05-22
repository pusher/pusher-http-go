package pusher

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
This is a JWT with header
{
  "alg": "HS256",
  "typ": "JWT"
}
and payload
{
  "exp": 1590155643,
  "iat": 1590155583,
  "iss": "278d425bdf160c739803",
  "sub": "callum",
  "user_info": {
    "foo": "bar"
  },
  "channels": [
    {
      "name": "private-foo"
    },
    {
      "name": "private-bar"
    },
    {
      "name": "presence-foobar"
    }
  ]
}signed by secret 7ad3773142a6692b25b8
*/
const expectedJWTWithUserID = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTAxNTU2NDMsImlhdCI6MTU5MDE1NTU4MywiaXNzIjoiMjc4ZDQyNWJkZjE2MGM3Mzk4MDMiLCJzdWIiOiJjYWxsdW0iLCJ1c2VyX2luZm8iOnsiZm9vIjoiYmFyIn0sImNoYW5uZWxzIjpbeyJuYW1lIjoicHJpdmF0ZS1mb28ifSx7Im5hbWUiOiJwcml2YXRlLWJhciJ9LHsibmFtZSI6InByZXNlbmNlLWZvb2JhciJ9XX0.SvBsw-QQtX8chLmhb3kjkMeXx-i28mO6EwspW_o-HDg"

// As above but without sub and user_info
const expectedJWTWithoutUserID = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTAxNTU2NDMsImlhdCI6MTU5MDE1NTU4MywiaXNzIjoiMjc4ZDQyNWJkZjE2MGM3Mzk4MDMiLCJjaGFubmVscyI6W3sibmFtZSI6InByaXZhdGUtZm9vIn0seyJuYW1lIjoicHJpdmF0ZS1iYXIifSx7Im5hbWUiOiJwcmVzZW5jZS1mb29iYXIifV19.PbGDBUr9IP5xQuvj7E4PI55nuIGnZQquKtSOlSaB-eo"

var now = time.Unix(1590155583, 0)

// This test is a little brittle because it relies on the order of the json keys
// being the same. Ideally we'd probably validate and parse the JWT instead.
func TestUserSessionAuthentication(t *testing.T) {
	client := setUpAuthClient()

	userID := "callum"
	userInfo := struct {
		Foo string `json:"foo"`
	}{"bar"}
	channelNames := []string{"private-foo", "private-bar", "presence-foobar"}

	expected := fmt.Sprintf(`{"auth":"%s"}`, expectedJWTWithUserID)
	result, err := client.authenticateSession(channelNames, userID, userInfo, now)
	assert.Equal(t, expected, string(result))
	assert.NoError(t, err)
}

// This test is a little brittle because it relies on the order of the json keys
// being the same. Ideally we'd probably validate and parse the JWT instead.
func TestSessionAuthentication(t *testing.T) {
	client := setUpAuthClient()

	channelNames := []string{"private-foo", "private-bar", "presence-foobar"}

	expected := fmt.Sprintf(`{"auth":"%s"}`, expectedJWTWithoutUserID)
	result, err := client.authenticateSession(channelNames, "", nil, now)
	assert.Equal(t, expected, string(result))
	assert.NoError(t, err)
}
