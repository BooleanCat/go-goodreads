package goodreads

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type Author struct {
	ID   string `xml:"id"`
	Name string `xml:"name"`
}

func (client Client) AuthorShow(id string) (Author, error) {
	type goodreadsResponse struct {
		Author Author `xml:"author"`
	}

	url := fmt.Sprintf("%s/author/show/%s.xml", goodreadsURL, id)
	response, err := client.doNewRequestWithKey(http.MethodGet, url, nil)
	if err != nil {
		return Author{}, err
	}
	defer closeIgnoreError(response.Body)

	if response.StatusCode != http.StatusOK {
		return Author{}, fmt.Errorf(`unexpected status code "%d"`, response.StatusCode)
	}

	var author goodreadsResponse
	if err := xml.NewDecoder(response.Body).Decode(&author); err != nil {
		return Author{}, fmt.Errorf("decode response: %w", err)
	}

	return author.Author, nil
}
