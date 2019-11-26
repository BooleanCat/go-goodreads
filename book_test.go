package goodreads_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/BooleanCat/go-goodreads"
	"github.com/BooleanCat/go-goodreads/fakes"
)

var _ = Describe("BookShow", func() {
	var (
		client       goodreads.Client
		fakeDoer     *fakes.FakeDoer
		responseBody *bytes.Buffer
	)

	BeforeEach(func() {
		responseBody = bytes.NewBufferString(`
			<goodreads_response>
				<book>
					<id>foo</id>
					<title>baz bar</title>
				</book>
			</goodreads_response>
		`)
		fakeDoer = new(fakes.FakeDoer)
		fakeDoer.DoReturns(&http.Response{
			Body:       ioutil.NopCloser(responseBody),
			StatusCode: http.StatusOK,
		}, nil)
		client = goodreads.NewClient(fakeDoer)
	})

	It("fetches the specified book by ID", func() {
		book, err := client.BookShow("foo")

		By("succeeding", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		By("decoding the book", func() {
			Expect(book).To(Equal(goodreads.Book{
				ID:    "foo",
				Title: "baz bar",
			}))
		})

		By("hitting the goodreads API correctly", func() {
			Expect(fakeDoer.DoCallCount()).To(Equal(1))
			request := fakeDoer.DoArgsForCall(0)
			Expect(request.Method).To(Equal(http.MethodGet))
			Expect(request.URL.String()).To(Equal("https://www.goodreads.com/book/show/foo.xml"))
		})
	})

	When("creating the request fails", func() {
		It("fails", func() {
			By("returning an error", func() {
				_, err := client.BookShow("%%%%%%")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("create request"))
			})

			By("not hitting the goodreads API", func() {
				Expect(fakeDoer.DoCallCount()).To(BeZero())
			})
		})
	})

	When("performing the request fails", func() {
		BeforeEach(func() {
			fakeDoer.DoReturns(nil, errors.New("oops"))
		})

		It("returns an error", func() {
			_, err := client.BookShow("foo")
			Expect(err).To(MatchError("do request: oops"))
		})
	})

	When("the response status code is not 200 OK", func() {
		BeforeEach(func() {
			fakeDoer.DoReturns(&http.Response{
				Body:       ioutil.NopCloser(responseBody),
				StatusCode: http.StatusMethodNotAllowed,
			}, nil)
		})

		It("returns an error", func() {
			_, err := client.BookShow("foo")
			Expect(err).To(MatchError(`unexpected status code "405"`))
		})
	})

	When("decoding the response body fails", func() {
		BeforeEach(func() {
			responseBody.Reset()
		})

		It("returns an error", func() {
			_, err := client.BookShow("foo")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("decode response"))
		})
	})
})
