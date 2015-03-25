package pusher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsingChannelsList(t *testing.T) {
	testJSON := []byte("{\"channels\":{\"presence-session-d41a439c438a100756f5-4bf35003e819bb138249-5cbTiUiPNGI\":{\"user_count\":1},\"presence-session-d41a439c438a100756f5-4bf35003e819bb138249-PbZ5E1pP8uF\":{\"user_count\":1},\"presence-session-d41a439c438a100756f5-4bf35003e819bb138249-oz6iqpSxMwG\":{\"user_count\":1}}}")
	expected := &ChannelsList{
		Channels: map[string]ChannelListItem{
			"presence-session-d41a439c438a100756f5-4bf35003e819bb138249-5cbTiUiPNGI": ChannelListItem{UserCount: 1},
			"presence-session-d41a439c438a100756f5-4bf35003e819bb138249-PbZ5E1pP8uF": ChannelListItem{UserCount: 1},
			"presence-session-d41a439c438a100756f5-4bf35003e819bb138249-oz6iqpSxMwG": ChannelListItem{UserCount: 1},
		},
	}
	result, err := unmarshalledChannelsList(testJSON)
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
}

func TestParsingChannel(t *testing.T) {
	testJSON := []byte("{\"user_count\":1,\"occupied\":true,\"subscription_count\":1}")
	channelName := "test"
	expected := &Channel{
		Name:              channelName,
		Occupied:          true,
		UserCount:         1,
		SubscriptionCount: 1,
	}
	result, err := unmarshalledChannel(testJSON, channelName)
	assert.Equal(t, expected, result)
	assert.NoError(t, err)

}

func TestParsingChannelUsers(t *testing.T) {
	testJSON := []byte("{\"users\":[{\"id\":\"red\"},{\"id\":\"blue\"}]}")
	expected := &Users{
		List: []User{User{Id: "red"}, User{Id: "blue"}},
	}
	result, err := unmarshalledChannelUsers(testJSON)
	assert.Equal(t, expected, result)
	assert.NoError(t, err)

}

func TestParsingTriggerResult(t *testing.T) {
	testJSON := []byte("{\"event_ids\":{\"test_channel\":\"eudhq17zrhfc4f\",\"another_channel\":\"eudhq17zrhfc74\"}}")
	expected := &BufferedEvents{
		EventIds: map[string]string{
			"test_channel":    "eudhq17zrhfc4f",
			"another_channel": "eudhq17zrhfc74",
		},
	}
	result, err := unmarshalledBufferedEvents(testJSON)
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
}

func TestParserError(t *testing.T) {
	testJSON := []byte("[];;[[p{{}}{{{}[][][]@£$@")
	_, err := unmarshalledChannelsList(testJSON)
	assert.Error(t, err)
}
