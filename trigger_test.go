package pusher

import (
	"github.com/pusher/pusher/requests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockDispatcher struct {
	mock.Mock
}

func (m *mockDispatcher) sendRequest(p *Pusher, request *requests.Request, params *requests.Params) (response []byte, err error) {
	args := m.Called(p, request, params)
	return args.Get(0).([]byte), args.Error(1)
}

func testTrigger(t *testing.T, req *requests.Request, triggerFunc func(*Pusher) (*TriggerResponse, error), expectJSON string) {
	mDispatcher := &mockDispatcher{}
	p := &Pusher{
		appID:      "id",
		key:        "key",
		secret:     "secret",
		dispatcher: mDispatcher,
	}

	expectedParams := &requests.Params{
		Body: []byte(expectJSON),
	}

	mDispatcher.
		On("sendRequest", p, req, expectedParams).
		Return([]byte("{}"), nil)

	_, err := triggerFunc(p)
	assert.NoError(t, err)

	mDispatcher.AssertExpectations(t)
}

func TestSimpleTrigger(t *testing.T) {
	testTrigger(t, requests.Trigger, func(p *Pusher) (*TriggerResponse, error) {
		return p.Trigger("test-channel", "my-event", map[string]string{"message": "hello world"})
	}, `{"name":"my-event","channels":["test-channel"],"data":"{\"message\":\"hello world\"}"}`)
}

func TestTriggerMulti(t *testing.T) {
	testTrigger(t, requests.Trigger, func(p *Pusher) (*TriggerResponse, error) {
		return p.TriggerMulti([]string{"test-channel-1", "test-channel-2"}, "my-event", "data")
	}, `{"name":"my-event","channels":["test-channel-1","test-channel-2"],"data":"data"}`)
}

func TestTriggerExclusive(t *testing.T) {
	testTrigger(t, requests.Trigger, func(p *Pusher) (*TriggerResponse, error) {
		return p.TriggerExclusive("test-channel", "my-event", "data", "123.12")
	}, `{"name":"my-event","channels":["test-channel"],"data":"data","socket_id":"123.12"}`)
}

func TestTriggerMultiExclusive(t *testing.T) {
	testTrigger(t, requests.Trigger, func(p *Pusher) (*TriggerResponse, error) {
		return p.TriggerMultiExclusive([]string{"channel-1", "channel-2"}, "event", "data", "123.12")
	}, `{"name":"event","channels":["channel-1","channel-2"],"data":"data","socket_id":"123.12"}`)
}

func TestTriggerBatch(t *testing.T) {
	testTrigger(t, requests.TriggerBatch, func(p *Pusher) (*TriggerResponse, error) {
		return p.TriggerBatch(
			[]Event{{
				Name:     "event-1",
				Channel:  "one",
				Data:     "data",
				SocketID: "123.12",
			}, {
				Name:    "event-2",
				Channel: "two",
				Data:    "data2",
			}},
		)
	}, `{"batch":[{"name":"event-1","channel":"one","data":"data","socket_id":"123.12"},{"name":"event-2","channel":"two","data":"data2"}]}`)
}
