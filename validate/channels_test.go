package validate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateSingleChannelTooLong(t *testing.T) {
	var longChannel string

	for i := 0; i <= 202; i++ {
		longChannel += "a"
	}
	err := Channels([]string{longChannel})
	expected :=
		`[pusher-http-go]: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa is over 200 characters`
	assert.EqualError(t, err, expected)
}

func TestValidateSingleChannelIllegalChars(t *testing.T) {
	badChannel := "∞¢#∞¢#∞¢§∞¢¢§#∞#"
	err := Channels([]string{badChannel})
	expected :=
		`[pusher-http-go]: ∞¢#∞¢#∞¢§∞¢¢§#∞# has illegal characters`
	assert.EqualError(t, err, expected)
}

func TestValidChannelsNoErrors(t *testing.T) {
	channels := []string{
		"test-channel",
		"yolo-channel",
	}
	err := Channels(channels)
	assert.NoError(t, err)
}

func TestMixOfErrorsAndNoErrors(t *testing.T) {
	channels := []string{
		"f¶§¶§∞∞¢",
		"hello",
		"gdfgfdogjfiodgjfodigjfdiogjfdiogjfidojgiofdjigoifdogjiofdjgojfdoigjiofdgjfdogijfdogifdjgiojfdgdfogjfdoijfdiogjfidogjiofdgjfodgijfdiogjfdoigjfdigodfgjoidgjreoigjreiogjeroigjeriogjireoireroigjergiorejigejrogierjoigjgroeigjoiegjoei",
		"waddup",
	}
	err := Channels(channels)
	expected :=
		`[pusher-http-go]: f¶§¶§∞∞¢ has illegal characters. gdfgfdogjfiodgjfodigjfdiogjfdiogjfidojgiofdjigoifdogjiofdjgojfdoigjiofdgjfdogijfdogifdjgiojfdgdfogjfdoijfdiogjfidogjiofdgjfodgijfdiogjfdoigjfdigodfgjoidgjreoigjreiogjeroigjeriogjireoireroigjergiorejigejrogierjoigjgroeigjoiegjoei is over 200 characters`
	assert.EqualError(t, err, expected)
}
