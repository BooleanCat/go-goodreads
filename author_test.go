package goodreads_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/BooleanCat/go-goodreads"
	"github.com/BooleanCat/go-goodreads/fakes"
	"github.com/BooleanCat/go-goodreads/httpclient"
)

func ExampleClient_AuthorShow() {
	client := goodreads.Client{
		Client: httpclient.RateLimited(http.DefaultClient, ticker),
	}

	book, err := client.AuthorShow("4764")
	if err != nil {
		panic(err)
	}

	fmt.Println(book.Name)
	// Output:
	// Philip K. Dick
}

func TestClient_AuthorShow(t *testing.T) {
	responseBody := bytes.NewBufferString(authorShowResponseBody)
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(responseBody),
		StatusCode: http.StatusOK,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	author, err := client.AuthorShow("foo")
	if err != nil {
		t.Fatalf(`expected error "%v" not to have occurred`, err)
	}

	want := goodreads.Author{ID: "foo", Name: "Baz"}
	if author != want {
		t.Fatalf(`expected author "%v" to equal "%v"`, author, want)
	}

	if fakeDoer.DoCallCount() != 1 {
		t.Fatal("expected request to have been performed")
	}

	request := fakeDoer.DoArgsForCall(0)

	if request.Method != http.MethodGet {
		t.Fatalf(`expected request method to equal "%s"`, http.MethodGet)
	}

	wantURL := "https://www.goodreads.com/author/show/foo.xml?key=key"
	if request.URL.String() != wantURL {
		t.Fatalf(`expected URL "%s" to equal %s`, request.URL.String(), wantURL)
	}
}

func TestClient_AuthorShow_CreateRequestFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	client := goodreads.Client{Client: nil, Key: "key"}

	_, err := client.AuthorShow("%%%%%%")
	if err == nil {
		t.Fatal("expected failure")
	}

	if !strings.Contains(err.Error(), "create request") {
		t.Fatalf(`expected error "%s" to contain "create request"`, err.Error())
	}

	if fakeDoer.DoCallCount() != 0 {
		t.Fatal("expected request not to have been performed")
	}
}

func TestClient_AuthorShow_DoRequestFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(nil, errors.New("oops"))
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.AuthorShow("foo")
	if err == nil {
		t.Fatal("expected failure")
	}

	if err.Error() != "do request: oops" {
		t.Fatalf(`expected error "%s" to equal "do request: oops"`, err.Error())
	}
}

func TestClient_AuthorShow_InvalidStatusCode(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusMethodNotAllowed,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.AuthorShow("foo")
	if err == nil {
		t.Fatal("expected failure")
	}

	want := `unexpected status code "405"`
	if err.Error() != want {
		t.Fatalf(`expected error "%s" to equal "%s`, err.Error(), want)
	}
}

func TestClient_AuthorShow_DecodeFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusOK,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.AuthorShow("foo")
	if err == nil {
		t.Fatal("expected failure")
	}

	if !strings.Contains(err.Error(), "decode response") {
		t.Fatalf(`expected error "%s" to contain "decode response"`, err.Error())
	}
}

const authorShowResponseBody string = `
	<goodreads_response>
		<author>
			<id>foo</id>
			<name>Baz</name>
		</author>
	</goodreads_response>
`
