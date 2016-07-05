package pusher

import (
	"github.com/pusher/pusher-http-go/requests"
	"net/http"
	"net/url"
)

type dispatcher interface {
	sendRequest(p *urlConfig, client *http.Client, request *requests.Request, params *requests.Params) (response []byte, err error)
}

type defaultDispatcher struct{}

func (d defaultDispatcher) sendRequest(p *urlConfig, client *http.Client, request *requests.Request, params *requests.Params) (response []byte, err error) {
	var u *url.URL
	if u, err = requestURL(p, request, params); err != nil {
		return
	}
	return request.Do(client, u, params.Body)
}
