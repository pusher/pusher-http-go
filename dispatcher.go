package pusher

import (
	"github.com/pusher/pusher-http-go/requests"
	"net/http"
)

type dispatcher interface {
	sendRequest(
		uc *urlConfig,
		client *http.Client,
		request *requests.Request,
		params *requests.Params,
	) (response []byte, err error)
}

type defaultDispatcher struct{}

func (d defaultDispatcher) sendRequest(
	uc *urlConfig,
	client *http.Client,
	request *requests.Request,
	params *requests.Params,
) ([]byte, error) {
	u, err := requestURL(uc, request, params)
	if err != nil {
		return nil, err
	}

	return request.Do(client, u, params.Body)
}
