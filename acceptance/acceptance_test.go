package acceptance_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/BooleanCat/go-goodreads"
	"github.com/BooleanCat/go-goodreads/httpclient"
)

var _ = Describe("go-goodreads", func() {
	var client goodreads.Client

	BeforeEach(func() {
		client = goodreads.NewClient(
			httpclient.WithKey(httpclient.RateLimited(
				http.DefaultClient, ticker,
			), goodreadsKey),
		)
	})

	Describe("/user/show/{user_id}.xml", func() {
		It("fetches the user", func() {
			user, err := client.UserShow("101333864")
			Expect(err).NotTo(HaveOccurred())
			Expect(user.UserName).To(Equal("tgodkin"))
		})
	})

	Describe("/author/show/{author_id}.xml", func() {
		It("fetches the author", func() {
			author, err := client.AuthorShow("4764")
			Expect(err).NotTo(HaveOccurred())
			Expect(author.Name).To(Equal("Philip K. Dick"))
		})
	})
})
