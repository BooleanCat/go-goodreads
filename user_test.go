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

var _ = Describe("UserShow", func() {
	var (
		client       goodreads.Client
		fakeDoer     *fakes.FakeDoer
		responseBody *bytes.Buffer
	)

	BeforeEach(func() {
		responseBody = bytes.NewBufferString(`
			<goodreads_response>
				<user>
					<id>foo</id>
					<name>Foo Bar</name>
					<user_name>fbar</user_name>
					<link><![CDATA[https://foo.com/fbar]]></link>
					<image_url><![CDATA[https://foo.com/fbar.png]]></image_url>
					<small_image_url><![CDATA[https://foo.com/fbarmini.png]]></small_image_url>
					<about>A test user</about>
					<age>30</age>
					<gender>male</gender>
					<location>London, The United Kingdom</location>
					<website>https://foo.com</website>
				</user>
			</goodreads_response>
		`)
		fakeDoer = new(fakes.FakeDoer)
		fakeDoer.DoReturns(&http.Response{
			Body:       ioutil.NopCloser(responseBody),
			StatusCode: http.StatusOK,
		}, nil)
		client = goodreads.NewClient(fakeDoer)
	})

	It("fetches the specified user by ID", func() {
		user, err := client.UserShow("foo")

		By("succeeding", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		By("decoding the user", func() {
			Expect(user).To(Equal(goodreads.User{
				ID:            "foo",
				Name:          "Foo Bar",
				UserName:      "fbar",
				Link:          "https://foo.com/fbar",
				ImageURL:      "https://foo.com/fbar.png",
				SmallImageURL: "https://foo.com/fbarmini.png",
				About:         "A test user",
				Age:           "30",
				Gender:        "male",
				Location:      "London, The United Kingdom",
				Website:       "https://foo.com",
			}))
		})

		By("hitting the goodreads API correctly", func() {
			Expect(fakeDoer.DoCallCount()).To(Equal(1))
			request := fakeDoer.DoArgsForCall(0)
			Expect(request.Method).To(Equal(http.MethodGet))
			Expect(request.URL.String()).To(Equal("https://www.goodreads.com/user/show/foo.xml"))
		})
	})

	When("creating the request fails", func() {
		It("fails", func() {
			By("returning an error", func() {
				_, err := client.UserShow("%%%%%%")
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
			_, err := client.UserShow("foo")
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
			_, err := client.UserShow("foo")
			Expect(err).To(MatchError(`unexpected status code "405"`))
		})
	})

	When("decoding the response body fails", func() {
		BeforeEach(func() {
			responseBody.Reset()
		})

		It("returns an error", func() {
			_, err := client.UserShow("foo")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("decode response"))
		})
	})
})
