package goodreads

import (
	"io"
	"net/http"
)

func (client Client) WithNewRequest(f func(string, string, io.Reader) (*http.Request, error)) Client {
	client.newRequestFunc = f
	return client
}
