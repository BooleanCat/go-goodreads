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
	"github.com/BooleanCat/go-goodreads/httputils"
)

func ExampleClient_AuthorShow() {
	transport := httputils.DripLimit(http.DefaultTransport, ticker)
	client := goodreads.Client{Client: &http.Client{Transport: transport}}

	book, err := client.AuthorShow(4764)
	if err != nil {
		panic(err)
	}

	fmt.Println(book.Name)
	// Output:
	// Philip K. Dick
}

func TestClient_AuthorShow(t *testing.T) {
	responseBody := bytes.NewBufferString(authorShowResponseBody)
	transport := new(fakes.FakeRoundTripper)
	transport.RoundTripReturns(&http.Response{
		Body:       ioutil.NopCloser(responseBody),
		StatusCode: http.StatusOK,
	}, nil)
	client := goodreads.Client{Client: &http.Client{Transport: transport}, Key: "key"}

	author, err := client.AuthorShow(123)
	assert.Nil(t, err)
	assert.Equal(t, author, goodreads.Author{
		ID:                   123,
		Name:                 "Baz",
		Link:                 "https://foo.com/author",
		FansCount:            42,
		AuthorFollowersCount: 50,
		LargeImageURL:        "https://foo.com/large.png",
		ImageURL:             "https://foo.com/image.png",
		SmallImageURL:        "https://foo.com/small.png",
		About:                "OK",
		Influences:           "bcat",
		WorksCount:           12,
		Gender:               "male",
		Hometown:             "London",
		BornAt:               "1945/12/03",
		DiedAt:               "1994/03/14",
		GoodreadsAuthor:      "baz",
		Books:                []goodreads.Book{{Title: "Mediocre Book"}},
	})

	assert.Equal(t, transport.RoundTripCallCount(), 1)
	request := transport.RoundTripArgsForCall(0)
	assert.Equal(t, request.Method, http.MethodGet)
	assert.Equal(t, request.URL.String(), "https://www.goodreads.com/author/show/123.xml?key=key")
}

func TestClient_AuthorShow_CreateRequestFails(t *testing.T) {
	transport := new(fakes.FakeRoundTripper)
	client := goodreads.Client{Client: nil, Key: "key", URL: "%%%"}

	_, err := client.AuthorShow(123)
	assert.ErrorMatches(t, err, `^create request: `)
	assert.Equal(t, transport.RoundTripCallCount(), 0)
}

func TestClient_AuthorShow_DoRequestFails(t *testing.T) {
	transport := new(fakes.FakeRoundTripper)
	transport.RoundTripReturns(nil, errors.New("oops"))
	client := goodreads.Client{Client: &http.Client{Transport: transport}, Key: "key"}

	_, err := client.AuthorShow(123)
	assert.ErrorMatches(t, err, `^do request: .*oops$`)
}

func TestClient_AuthorShow_InvalidStatusCode(t *testing.T) {
	transport := new(fakes.FakeRoundTripper)
	transport.RoundTripReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusMethodNotAllowed,
	}, nil)
	client := goodreads.Client{Client: &http.Client{Transport: transport}, Key: "key"}

	_, err := client.AuthorShow(123)
	assert.ErrorMatches(t, err, `^unexpected status code "405"$`)
}

func TestClient_AuthorShow_DecodeFails(t *testing.T) {
	transport := new(fakes.FakeRoundTripper)
	transport.RoundTripReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusOK,
	}, nil)
	client := goodreads.Client{Client: &http.Client{Transport: transport}, Key: "key"}

	_, err := client.AuthorShow(123)
	assert.ErrorMatches(t, err, `^decode response: `)
}

const authorShowResponseBody string = `
	<goodreads_response>
		<author>
			<id>123</id>
			<name>Baz</name>
			<link>https://foo.com/author</link>
			<fans_count>42</fans_count>
			<author_followers_count>50</author_followers_count>
			<large_image_url>https://foo.com/large.png</large_image_url>
			<image_url>https://foo.com/image.png</image_url>
			<small_image_url>https://foo.com/small.png</small_image_url>
			<about>OK</about>
			<influences>bcat</influences>
			<works_count>12</works_count>
			<gender>male</gender>
			<hometown>London</hometown>
			<born_at>1945/12/03</born_at>
			<died_at>1994/03/14</died_at>
			<goodreads_author>baz</goodreads_author>
			<books>
				<book>
					<title>Mediocre Book</title>
				</book>
			</books>
		</author>
	</goodreads_response>
`
