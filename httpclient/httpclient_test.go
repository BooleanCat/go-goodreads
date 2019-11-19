package httpclient_test

import (
	"errors"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/BooleanCat/go-goodreads/httpclient"
	"github.com/BooleanCat/go-goodreads/httpclient/httpclientfakes"
)

var _ = Describe("RateLimitedClient", func() {
	var fakeDoer *httpclientfakes.FakeDoer

	BeforeEach(func() {
		fakeDoer = new(httpclientfakes.FakeDoer)
	})

	It("makes http requests", func() {
		ticker := time.NewTicker(time.Millisecond)
		defer ticker.Stop()
		client := httpclient.RateLimited(fakeDoer, ticker)

		_, err := client.Do(new(http.Request))

		By("succeeding", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		By("invoking the delegate do", func() {
			Expect(fakeDoer.DoCallCount()).To(Equal(1))
		})
	})

	When("the delegate do fails", func() {
		BeforeEach(func() {
			fakeDoer.DoReturns(nil, errors.New("oops"))
		})

		It("returns the error", func() {
			ticker := time.NewTicker(time.Millisecond)
			defer ticker.Stop()
			client := httpclient.RateLimited(fakeDoer, ticker)

			_, err := client.Do(new(http.Request))
			Expect(err).To(MatchError("oops"))
		})
	})
})
