package goodreads_test

import (
	"os"
	"testing"
	"time"
)

var ticker *time.Ticker

func TestMain(m *testing.M) {
	ticker = time.NewTicker(time.Second * 2)

	exitCode := m.Run()
	ticker.Stop()
	os.Exit(exitCode)
}
