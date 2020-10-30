package goodreads

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/BooleanCat/go-goodreads/param"
)

// Client is used to communicate with the Goodreads API. A zero value client
// is valid for use, however most API methods will require one or both of Key
// and Secret to be set.
//
// If the GOODREADS_KEY environmental variable is set, it will be used in the
// case of client.Key being an empty string.
type Client struct {
	Client *http.Client
	URL    string
	Key    string
	Secret string
}

func (client Client) String() string {
	return fmt.Sprintf("{%v}", client.Client)
}

func (client Client) goodreadsKey() (string, error) {
	if client.Key != "" {
		return client.Key, nil
	}

	if key := os.Getenv("GOODREADS_KEY"); key != "" {
		return key, nil
	}

	return "", ErrAPIKeyNotSet{}
}

const defaultGoodreadsURL = "https://www.goodreads.com"

func closeIgnoreError(c io.Closer) {
	_ = c.Close()
}

func (client Client) newRequestWithKey(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	key, err := client.goodreadsKey()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	return param.Apply(request, param.APIKey(key)), nil
}

func (client Client) getClient() *http.Client {
	if client.Client == nil {
		return http.DefaultClient
	}

	return client.Client
}

func (client Client) getURL() string {
	if client.URL == "" {
		return defaultGoodreadsURL
	}

	return client.URL
}

//go:generate counterfeiter --generate
//counterfeiter:generate -o internal/fakes/roundtripper.go net/http.RoundTripper
