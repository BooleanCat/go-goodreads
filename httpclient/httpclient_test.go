package httpclient_test

import (
	"errors"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/BooleanCat/go-goodreads/fakes"
	"github.com/BooleanCat/go-goodreads/httpclient"
)

var _ = Describe("RateLimitedClient", func() {
	var fakeDoer *fakes.FakeDoer

	BeforeEach(func() {
		fakeDoer = new(fakes.FakeDoer)
	})

	Describe("RateLimitedClient", func() {
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

	Describe("KeyClient", func() {
		It(`add the "key" query param`, func() {
			request, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
			Expect(err).NotTo(HaveOccurred())
			client := httpclient.WithKey(fakeDoer, "foo")
			_, err = client.Do(request)

			By("succeeding", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			By("invoking the delegate do", func() {
				Expect(fakeDoer.DoCallCount()).To(Equal(1))
			})

			By("adding the query param", func() {
				got := fakeDoer.DoArgsForCall(0)
				Expect(got.URL.String()).To(ContainSubstring("key=foo"))
			})
		})
	})
})
