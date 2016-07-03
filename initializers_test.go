package pusher

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitialisationFromURL(t *testing.T) {
	url := "http://feaf18a411d3cb9216ee:fec81108d90e1898e17a@api.pusherapp.com/apps/104060"
	client, _ := NewFromURL(url)
	expectedClient := NewWithOptions("104060", "feaf18a411d3cb9216ee", "fec81108d90e1898e17a", Options{Host: "api.pusherapp.com"})
	assert.ObjectsAreEqual(expectedClient, client)
}

func TestURLInitErrorNoSecret(t *testing.T) {
	url := "http://fec81108d90e1898e17a@api.pusherapp.com/apps"
	client, err := NewFromURL(url)
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestURLInitHTTPS(t *testing.T) {
	url := "https://key:secret@api.pusherapp.com/apps/104060"
	client, _ := NewFromURL(url)
	assert.True(t, client.(*Pusher).Secure)
}

func TestURLInitErrorNoID(t *testing.T) {
	url := "http://fec81108d90e1898e17a@api.pusherapp.com/apps"
	client, err := NewFromURL(url)
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestInitialisationFromENV(t *testing.T) {
	os.Setenv("PUSHER_URL", "http://feaf18a411d3cb9216ee:fec81108d90e1898e17a@api.pusherapp.com/apps/104060")
	client, _ := NewFromEnv("PUSHER_URL")
	expectedClient := NewWithOptions("104060", "feaf18a411d3cb9216ee", "fec81108d90e1898e17a", Options{Host: "api.pusherapp.com"})
	assert.Equal(t, expectedClient, client)
}
