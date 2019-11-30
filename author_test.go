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
	client := goodreads.Client{
		Client: httputils.DripLimit(http.DefaultClient, ticker),
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
	assert.Nil(t, err)
	assert.Equal(t, author, goodreads.Author{
		ID:                   "foo",
		Name:                 "Baz",
		Link:                 "https://foo.com/author",
		FansCount:            42,
		AuthorFollowersCount: 50,
		LargeImageURL:        "https://foo.com/large.png",
		ImageURL:             "https://foo.com/image.png",
		SmallImageURL:        "https://foo.com/small.png",
		About:                "OK",
		Influences:           "bcat",
		WorksCount:           "12",
		Gender:               "male",
		Hometown:             "London",
		BornAt:               "1945/12/03",
		DiedAt:               "1994/03/14",
		GoodreadsAuthor:      "baz",
		Books: goodreads.Books{Books: []goodreads.Book{
			{Title: "Mediocre Book"},
		}},
	})

	assert.Equal(t, fakeDoer.DoCallCount(), 1)
	request := fakeDoer.DoArgsForCall(0)
	assert.Equal(t, request.Method, http.MethodGet)
	assert.Equal(t, request.URL.String(), "https://www.goodreads.com/author/show/foo.xml?key=key")
}

func TestClient_AuthorShow_CreateRequestFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	client := goodreads.Client{Client: nil, Key: "key"}

	_, err := client.AuthorShow("%%%%%%")
	assert.ErrorMatches(t, err, `^create request: `)
	assert.Equal(t, fakeDoer.DoCallCount(), 0)
}

func TestClient_AuthorShow_DoRequestFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(nil, errors.New("oops"))
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.AuthorShow("foo")
	assert.ErrorMatches(t, err, `^do request: oops$`)
}

func TestClient_AuthorShow_InvalidStatusCode(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusMethodNotAllowed,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.AuthorShow("foo")
	assert.ErrorMatches(t, err, `^unexpected status code "405"$`)
}

func TestClient_AuthorShow_DecodeFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusOK,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.AuthorShow("foo")
	assert.ErrorMatches(t, err, `^decode response: `)
}

const authorShowResponseBody string = `
	<goodreads_response>
		<author>
			<id>foo</id>
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
