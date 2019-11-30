package goodreads

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type Book struct {
	ID    string `xml:"id"`
	Title string `xml:"title"`
}

func (client Client) BookShow(id string) (Book, error) {
	type goodreadsResponse struct {
		Book Book `xml:"book"`
	}

	url := fmt.Sprintf("%s/book/show/%s.xml", goodreadsURL, id)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Book{}, fmt.Errorf("create request: %w", err)
	}

	request, err = client.addGoodreadsKeyQueryParam(request)
	if err != nil {
		return Book{}, err
	}

	response, err := client.Client.Do(request)
	if err != nil {
		return Book{}, fmt.Errorf("do request: %w", err)
	}
	defer closeIgnoreError(response.Body)

	if response.StatusCode != http.StatusOK {
		return Book{}, fmt.Errorf(`unexpected status code "%d"`, response.StatusCode)
	}

	var book goodreadsResponse
	if err := xml.NewDecoder(response.Body).Decode(&book); err != nil {
		return Book{}, fmt.Errorf("decode response: %w", err)
	}

	return book.Book, nil
}
