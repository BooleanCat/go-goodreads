package goodreads

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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

// OptionKey adds the goodreads API key to API calls.
func OptionKey(key string) func(url.Values) url.Values {
	return func(values url.Values) url.Values {
		values.Set("key", key)

		return values
	}
}

type option func(url.Values) url.Values

func setOptions(request *http.Request, options ...option) *http.Request {
	query := request.URL.Query()

	for _, o := range options {
		o(query)
	}

	request.URL.RawQuery = query.Encode()

	return request
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

	return setOptions(request, OptionKey(key)), nil
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
