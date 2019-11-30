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

func ExampleClient_BookShow() {
	client := goodreads.Client{
		Client: httpclient.RateLimited(http.DefaultClient, ticker),
	}

	book, err := client.BookShow("36402034")
	if err != nil {
		panic(err)
	}

	fmt.Println(book.Title)
	// Output:
	// Do Androids Dream of Electric Sheep?
}

func TestClient_BookShow(t *testing.T) {
	responseBody := bytes.NewBufferString(bookShowResponseBody)
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(responseBody),
		StatusCode: http.StatusOK,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	book, err := client.BookShow("foo")
	if err != nil {
		t.Fatalf(`expected error "%v" not to have occurred`, err)
	}

	want := goodreads.Book{ID: "foo", Title: "baz bar"}
	if book != want {
		t.Fatalf(`expected book "%v" to equal "%v"`, book, want)
	}

	if fakeDoer.DoCallCount() != 1 {
		t.Fatal("expected request to have been performed")
	}

	request := fakeDoer.DoArgsForCall(0)

	if request.Method != http.MethodGet {
		t.Fatalf(`expected request method to equal "%s"`, http.MethodGet)
	}

	wantURL := "https://www.goodreads.com/book/show/foo.xml?key=key"
	if request.URL.String() != wantURL {
		t.Fatalf(`expected URL "%s" to equal %s`, request.URL.String(), wantURL)
	}
}

func TestClient_BookShow_CreateRequestFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusOK,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.BookShow("%%%%%%")
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

func TestClient_BookShow_DoRequestFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(nil, errors.New("oops"))
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.BookShow("foo")
	if err == nil {
		t.Fatal("expected failure")
	}

	if err.Error() != "do request: oops" {
		t.Fatalf(`expected error "%s" to equal "do request: oops"`, err.Error())
	}
}

func TestClient_BookShow_InvalidStatusCode(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusMethodNotAllowed,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.BookShow("foo")
	if err == nil {
		t.Fatal("expected failure")
	}

	want := `unexpected status code "405"`
	if err.Error() != want {
		t.Fatalf(`expected error "%s" to equal "%s`, err.Error(), want)
	}
}

func TestClient_BookShow_DecodeFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusOK,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.BookShow("foo")
	if err == nil {
		t.Fatal("expected failure")
	}

	if !strings.Contains(err.Error(), "decode response") {
		t.Fatalf(`expected error "%s" to contain "decode response"`, err.Error())
	}
}

const bookShowResponseBody string = `
	<goodreads_response>
		<book>
			<id>foo</id>
			<title>baz bar</title>
		</book>
	</goodreads_response>
`
