package httputils

import (
	"net/http"
	"time"
)

type doer interface {
	Do(r *http.Request) (*http.Response, error)
}

type DripLimitClient struct {
	delegate doer
	ticker   *time.Ticker
}

func DripLimit(delegate doer, ticker *time.Ticker) DripLimitClient {
	return DripLimitClient{
		delegate: delegate,
		ticker:   ticker,
	}
}

func (client DripLimitClient) Do(request *http.Request) (*http.Response, error) {
	<-client.ticker.C
	return client.delegate.Do(request)
}

var _ doer = DripLimitClient{}
