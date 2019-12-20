package goodreads_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/BooleanCat/go-goodreads"
	"github.com/BooleanCat/go-goodreads/assert"
)

var ticker *time.Ticker

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
