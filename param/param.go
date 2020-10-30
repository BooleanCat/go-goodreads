package param

import (
	"fmt"
	"net/http"
	"net/url"
)

// Param is a mutation of a URL's values.
type Param func(url.Values) url.Values

// Apply given params to request.
func Apply(request *http.Request, params ...Param) *http.Request {
	query := request.URL.Query()

	for _, p := range params {
		p(query)
	}

	request.URL.RawQuery = query.Encode()

	return request
}

// APIKey adds the goodreads API key to API calls.
func APIKey(key string) Param {
	return func(values url.Values) url.Values {
		values.Set("key", key)

		return values
	}
}

var _ Param = APIKey("")

// TextOnly shows only reviews with text.
func TextOnly(values url.Values) url.Values {
	values.Set("text_only", "true")

	return values
}

var _ Param = TextOnly

// Rating shows only reviews with a given rating.
func Rating(r float32) func(url.Values) url.Values {
	return func(values url.Values) url.Values {
		values.Set("rating", fmt.Sprintf("%.2f", r))

		return values
	}
}

var _ Param = Rating(0)
