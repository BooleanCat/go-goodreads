package httputils

import (
	"net/http"
	"time"
)

type DripLimitTransport struct {
	delegate http.RoundTripper
	ticker   *time.Ticker
}

func DripLimit(delegate http.RoundTripper, ticker *time.Ticker) DripLimitTransport {
	return DripLimitTransport{
		delegate: delegate,
		ticker:   ticker,
	}
}

func (client DripLimitTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	<-client.ticker.C
	return client.delegate.RoundTrip(request)
}

var _ http.RoundTripper = DripLimitTransport{}
