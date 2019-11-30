package httpclient_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/BooleanCat/go-goodreads/assert"
	"github.com/BooleanCat/go-goodreads/fakes"
	"github.com/BooleanCat/go-goodreads/httpclient"
)

func TestRateLimited(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	client := httpclient.RateLimited(fakeDoer, ticker)

	_, err := client.Do(new(http.Request))
	assert.Nil(t, err)
	assert.Equal(t, fakeDoer.DoCallCount(), 1)
}

func TestRateLimited_DelegateDoFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(nil, errors.New("oops"))

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	client := httpclient.RateLimited(fakeDoer, ticker)

	_, err := client.Do(new(http.Request))
	assert.ErrorMatches(t, err, `oops`)
}
