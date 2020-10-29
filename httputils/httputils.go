package httputils

import (
	"net/http"
	"time"
)

// DripLimitTransport rate limits requests by regular intervals.
type DripLimitTransport struct {
	delegate http.RoundTripper
	ticker   *time.Ticker
}

// DripLimit creates a new DripLimitTransport. The caller must stop the time.Ticker.
func DripLimit(delegate http.RoundTripper, ticker *time.Ticker) DripLimitTransport {
	return DripLimitTransport{
		delegate: delegate,
		ticker:   ticker,
	}
}

// RoundTrip implements http.RoundTripper.
func (client DripLimitTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	<-client.ticker.C

	return client.delegate.RoundTrip(request)
}

var _ http.RoundTripper = DripLimitTransport{}
