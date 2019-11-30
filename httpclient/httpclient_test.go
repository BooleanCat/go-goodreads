package httpclient_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/BooleanCat/go-goodreads/fakes"
	"github.com/BooleanCat/go-goodreads/httpclient"
)

func TestRateLimited(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	client := httpclient.RateLimited(fakeDoer, ticker)

	_, err := client.Do(new(http.Request))
	if err != nil {
		t.Fatalf(`expected error "%v" not to have occurred`, err)
	}

	if fakeDoer.DoCallCount() != 1 {
		t.Fatal("expected request to have been performed")
	}
}

func TestRateLimited_DelegateDoFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(nil, errors.New("oops"))

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	client := httpclient.RateLimited(fakeDoer, ticker)

	_, err := client.Do(new(http.Request))
	if err == nil {
		t.Fatal("expected failure")
	}

	if err.Error() != "oops" {
		t.Fatalf(`expected error "%s" to equal "oops"`, err.Error())
	}
}
