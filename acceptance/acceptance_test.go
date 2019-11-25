package acceptance_test

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/BooleanCat/go-goodreads"
	"github.com/BooleanCat/go-goodreads/httpclient"
)

var _ = Describe("go-goodreads", func() {
	var (
		ticker *time.Ticker
		client goodreads.Client
	)

	BeforeEach(func() {
		ticker = time.NewTicker(time.Second * 2)
		client = goodreads.NewClient(
			httpclient.WithKey(httpclient.RateLimited(
				http.DefaultClient, ticker,
			), goodreadsKey),
		)
	})

	AfterEach(func() {
		ticker.Stop()
	})

	Describe("/user/show/{user_id}.xml", func() {
		It("fetches the user", func() {
			user, err := client.UserShow("101333864")
			Expect(err).NotTo(HaveOccurred())
			Expect(user.UserName).To(Equal("tgodkin"))
		})
	})
})
