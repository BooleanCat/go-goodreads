package goodreads

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Client struct {
	Client doer
	Key    string
	Secret string

	newRequestFunc newRequestFunc
}

func (client Client) String() string {
	return fmt.Sprintf("{%v}", client.Client)
}

func (client Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	if client.newRequestFunc != nil {
		return client.newRequestFunc(method, url, body)
	}

	return http.NewRequest(method, url, body)
}

func (client Client) addGoodreadsKeyQueryParam(request *http.Request) (*http.Request, error) {
	key, err := client.goodreadsKey()
	if err != nil {
		return nil, err
	}

	query := request.URL.Query()
	query.Add("key", key)
	request.URL.RawQuery = query.Encode()

	return request, nil
}

func (client Client) goodreadsKey() (string, error) {
	if client.Key != "" {
		return client.Key, nil
	}

	if key := os.Getenv("GOODREADS_KEY"); key != "" {
		return key, nil
	}

	return "", errors.New("goodreads API key not set")
}

const goodreadsURL = "https://www.goodreads.com"

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

type newRequestFunc func(string, string, io.Reader) (*http.Request, error)

func closeIgnoreError(c io.Closer) {
	_ = c.Close()
}

func (client Client) doNewRequestWithKey(method, url string, body io.Reader) (*http.Response, error) {
	request, err := client.newRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	request, err = client.addGoodreadsKeyQueryParam(request)
	if err != nil {
		return nil, err
	}

	response, err := client.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return response, nil
}

//go:generate counterfeiter --generate
//counterfeiter:generate -o fakes/doer.go . doer
//counterfeiter:generate -o fakes/request.go . newRequestFunc
