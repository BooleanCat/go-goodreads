package acceptance_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "acceptance suite")
}

var goodreadsKey string

var _ = BeforeSuite(func() {
	goodreadsKey = os.Getenv("GOODREADS_KEY")
	if goodreadsKey == "" {
		Fail("GOODREADS_KEY must be set")
	}
})
