package acceptance_test

import (
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "acceptance suite")
}

var (
	goodreadsKey string
	ticker       *time.Ticker
)

var _ = BeforeSuite(func() {
	ticker = time.NewTicker(time.Second * 2)
	goodreadsKey = os.Getenv("GOODREADS_KEY")
	if goodreadsKey == "" {
		Fail("GOODREADS_KEY must be set")
	}
})

var _ = AfterSuite(func() {
	ticker.Stop()
})
