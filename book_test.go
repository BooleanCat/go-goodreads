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

func ExampleClient_BookShow() {
	client := goodreads.Client{
		Client: httputils.DripLimit(http.DefaultClient, ticker),
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
	assert.Equal(t, book, goodreads.Book{
		ID:                 "foo",
		Title:              "baz bar",
		ISBN:               "isbn",
		ISBN13:             "isbn13",
		ASIN:               "asin",
		KindleASIN:         "kindle asin",
		MarketplaceID:      "foobar",
		CountryCode:        "GB",
		ImageURL:           "https://foo.com/bar.png",
		SmallImageURL:      "https://foo.com/baz.png",
		PublicationYear:    "2019",
		PublicationMonth:   "2",
		PublicationDay:     "22",
		Publisher:          "bcat",
		LanguageCode:       "eng",
		IsEbook:            "true",
		Description:        "What a book.",
		AverageRating:      "4.09",
		NumPages:           "201",
		Format:             "Kindle",
		EditionInformation: "Best edition",
		RatingsCount:       "98",
		TextReviewsCount:   "42",
		URL:                "https://foo.com/book",
		Link:               "https://bar.com/book",
		Work: goodreads.Work{
			ID:                             42,
			BooksCount:                     5,
			BestBookID:                     765,
			ReviewsCount:                   653,
			RatingsSum:                     1000,
			RatingsCount:                   400,
			TextReviewsCount:               50,
			OriginalPublicationYear:        2019,
			OriginalPublicationMonth:       4,
			OriginalPublicationDay:         30,
			OriginalTitle:                  "Bar",
			OriginalLanguageID:             0,
			MediaType:                      "book",
			RatingDist:                     "5:110406|4:126699|3:56731|2:10593|1:2539|total:306968",
			DescUserID:                     788,
			DefaultChapteringBookID:        14,
			DefaultDescriptionLanguageCode: "eng",
			WorkURI:                        "https://foo.com",
		},
		Authors: goodreads.Authors{Authors: []goodreads.Author{
			{Name: "bcat"},
		}},
	})

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
			<isbn>isbn</isbn>
			<isbn13>isbn13</isbn13>
			<asin>asin</asin>
			<kindle_asin>kindle asin</kindle_asin>
			<marketplace_id>foobar</marketplace_id>
			<country_code>GB</country_code>
			<image_url>https://foo.com/bar.png</image_url>
			<small_image_url>https://foo.com/baz.png</small_image_url>
			<publication_year>2019</publication_year>
			<publication_month>2</publication_month>
			<publication_day>22</publication_day>
			<publisher>bcat</publisher>
			<language_code>eng</language_code>
			<is_ebook>true</is_ebook>
			<description>What a book.</description>
			<average_rating>4.09</average_rating>
			<num_pages>201</num_pages>
			<format>Kindle</format>
			<edition_information>Best edition</edition_information>
			<ratings_count>98</ratings_count>
			<text_reviews_count>42</text_reviews_count>
			<url>https://foo.com/book</url>
			<link>https://bar.com/book</link>
			<work>
				<id>42</id>
				<books_count>5</books_count>
				<best_book_id>765</best_book_id>
				<reviews_count>653</reviews_count>
				<ratings_sum>1000</ratings_sum>
				<ratings_count>400</ratings_count>
				<text_reviews_count>50</text_reviews_count>
				<original_publication_year>2019</original_publication_year>
				<original_publication_month>4</original_publication_month>
				<original_publication_day>30</original_publication_day>
				<original_title>Bar</original_title>
				<original_language_id />
				<media_type>book</media_type>
				<rating_dist>5:110406|4:126699|3:56731|2:10593|1:2539|total:306968</rating_dist>
				<desc_user_id>788</desc_user_id>
				<default_chaptering_book_id>14</default_chaptering_book_id>
				<default_description_language_code>eng</default_description_language_code>
				<work_uri>https://foo.com</work_uri>
			</work>
			<authors>
				<author>
					<name>bcat</name>
				</author>
			</authors>
		</book>
	</goodreads_response>
`
