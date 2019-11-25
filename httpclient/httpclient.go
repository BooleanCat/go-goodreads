package httpclient

import (
	"net/http"
	"time"
)

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

type KeyClient struct {
	delegate doer
	key      string
}

func WithKey(client doer, key string) KeyClient {
	return KeyClient{
		delegate: client,
		key:      key,
	}
}

func (client KeyClient) Do(request *http.Request) (*http.Response, error) {
	query := request.URL.Query()
	query.Add("key", client.key)
	request.URL.RawQuery = query.Encode()
	return client.delegate.Do(request)
}
