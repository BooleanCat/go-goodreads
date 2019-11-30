package httputils_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/BooleanCat/go-goodreads/assert"
	"github.com/BooleanCat/go-goodreads/fakes"
	"github.com/BooleanCat/go-goodreads/httputils"
)

func TestDripLimit(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	client := httputils.DripLimit(fakeDoer, ticker)

	_, err := client.Do(new(http.Request))
	assert.Nil(t, err)
	assert.Equal(t, fakeDoer.DoCallCount(), 1)
}

func TestDripLimit_DelegateDoFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(nil, errors.New("oops"))

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	client := httputils.DripLimit(fakeDoer, ticker)

	_, err := client.Do(new(http.Request))
	assert.ErrorMatches(t, err, `oops`)
}
