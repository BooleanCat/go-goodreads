package httputils_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/BooleanCat/go-goodreads/httputils"
	"github.com/BooleanCat/go-goodreads/internal/assert"
	"github.com/BooleanCat/go-goodreads/internal/fakes"
)

func TestDripLimit(t *testing.T) {
	transport := new(fakes.FakeRoundTripper)

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	client := httputils.DripLimit(transport, ticker)

	_, err := client.RoundTrip(new(http.Request))
	assert.Nil(t, err)
	assert.Equal(t, transport.RoundTripCallCount(), 1)
}

func TestDripLimit_DelegateDoFails(t *testing.T) {
	transport := new(fakes.FakeRoundTripper)
	transport.RoundTripReturns(nil, errors.New("oops"))

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	client := httputils.DripLimit(transport, ticker)

	_, err := client.RoundTrip(new(http.Request))
	assert.ErrorMatches(t, err, `oops`)
}
