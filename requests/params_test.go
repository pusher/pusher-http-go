package requests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockClock struct{}

func (m mockClock) Now() string {
	return "123456789"
}

func init() {
	authClock = mockClock{}
}

func TestParamURLValuesNoBodyOrQueries(t *testing.T) {
	params := &Params{}
	values := params.URLValues("key")
	actual := values.Encode()
	expected := "auth_key=key&auth_timestamp=123456789&auth_version=1.0"
	assert.Equal(t, expected, actual)
}

func TestParamURLValuesWithBodyMD5(t *testing.T) {
	params := &Params{
		Body: []byte("yolo123"),
	}
	values := params.URLValues("key2")
	actual := values.Encode()
	expected := "auth_key=key2&auth_timestamp=123456789&auth_version=1.0&body_md5=b48fee99c626f0634db4bf8f5d2d54b2"
	assert.Equal(t, expected, actual)
}

func TestAdditionalQueriesAddedToString(t *testing.T) {
	params := &Params{
		Body: []byte("yolo123"),
		Queries: map[string]string{
			"q": "123",
		},
	}
	values := params.URLValues("key3")
	actual := values.Encode()
	expected := "auth_key=key3&auth_timestamp=123456789&auth_version=1.0&body_md5=b48fee99c626f0634db4bf8f5d2d54b2&q=123"
	assert.Equal(t, expected, actual)
}
