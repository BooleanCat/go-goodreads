package goodreads_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/BooleanCat/go-goodreads"
	"github.com/BooleanCat/go-goodreads/internal/assert"
)

var ticker *time.Ticker //nolint:gochecknoglobals

func TestMain(m *testing.M) {
	ticker = time.NewTicker(time.Second * 2)

	exitCode := m.Run()

	ticker.Stop()
	os.Exit(exitCode)
}

func TestClient_String(t *testing.T) {
	client := goodreads.Client{Key: "foo", Secret: "bar"}
	assert.DoesNotContainSubstring(t, fmt.Sprint(client), "foo")
	assert.DoesNotContainSubstring(t, fmt.Sprint(client), "bar")
}

type fakeErr struct{}

func (err fakeErr) Error() string {
	return "oops"
}

var _ error = fakeErr{}
