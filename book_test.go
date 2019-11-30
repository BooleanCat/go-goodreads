package goodreads_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/BooleanCat/go-goodreads"
	"github.com/BooleanCat/go-goodreads/assert"
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
	assert.Nil(t, err)
	assert.Equal(t, book, goodreads.Book{ID: "foo", Title: "baz bar"})

	assert.Equal(t, fakeDoer.DoCallCount(), 1)
	request := fakeDoer.DoArgsForCall(0)
	assert.Equal(t, request.Method, http.MethodGet)
	assert.Equal(t, request.URL.String(), "https://www.goodreads.com/book/show/foo.xml?key=key")
}

func TestClient_BookShow_CreateRequestFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusOK,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.BookShow("%%%%%%")
	assert.ErrorMatches(t, err, `^create request: `)
	assert.Equal(t, fakeDoer.DoCallCount(), 0)
}

func TestClient_BookShow_DoRequestFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(nil, errors.New("oops"))
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.BookShow("foo")
	assert.ErrorMatches(t, err, `^do request: oops$`)
}

func TestClient_BookShow_InvalidStatusCode(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusMethodNotAllowed,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.BookShow("foo")
	assert.ErrorMatches(t, err, `^unexpected status code "405"$`)
}

func TestClient_BookShow_DecodeFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusOK,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.BookShow("foo")
	assert.ErrorMatches(t, err, `^decode response: `)
}

const bookShowResponseBody string = `
	<goodreads_response>
		<book>
			<id>foo</id>
			<title>baz bar</title>
		</book>
	</goodreads_response>
`
