package httpclient

import (
	"net/http"
	"time"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . doer
type doer interface {
	Do(r *http.Request) (*http.Response, error)
}

type RateLimitedClient struct {
	delegate doer
	ticker   *time.Ticker
}

func RateLimited(delegate doer, ticker *time.Ticker) RateLimitedClient {
	return RateLimitedClient{
		delegate: delegate,
		ticker:   ticker,
	}
}

func (client RateLimitedClient) Do(request *http.Request) (*http.Response, error) {
	<-client.ticker.C
	return client.delegate.Do(request)
}

var _ doer = RateLimitedClient{}
