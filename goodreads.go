package goodreads

import (
	"io"
	"net/http"
)

type Client struct {
	client doer
}

func NewClient(delegate doer) Client {
	return Client{client: delegate}
}

const goodreadsURL = "https://www.goodreads.com"

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

func closeIgnoreError(c io.Closer) {
	_ = c.Close()
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/doer.go . doer
