package pusher

import (
	"github.com/pusher/pusher/requests"
	"net/url"
)

type dispatcher interface {
	sendRequest(p *Pusher, request *requests.Request, params *requests.Params) (response []byte, err error)
}

type defaultDispatcher struct{}

func (d defaultDispatcher) sendRequest(p *Pusher, request *requests.Request, params *requests.Params) (response []byte, err error) {
	var u *url.URL
	if u, err = requestURL(p, request, params); err != nil {
		return
	}
	return request.Do(p.httpClient(), u, params.Body)
}
