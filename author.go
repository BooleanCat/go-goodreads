package goodreads

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

// An Author contains information about an author as defined by Goodreads.
type Author struct {
	ID                   int    `xml:"id"`
	Name                 string `xml:"name"`
	Link                 string `xml:"link"`
	FansCount            int    `xml:"fans_count"`
	AuthorFollowersCount int    `xml:"author_followers_count"`
	LargeImageURL        string `xml:"large_image_url"`
	ImageURL             string `xml:"image_url"`
	SmallImageURL        string `xml:"small_image_url"`
	About                string `xml:"about"`
	Influences           string `xml:"influences"`
	WorksCount           int    `xml:"works_count"`
	Gender               string `xml:"gender"`
	Hometown             string `xml:"hometown"`
	BornAt               string `xml:"born_at"`
	DiedAt               string `xml:"died_at"`
	GoodreadsAuthor      string `xml:"goodreads_author"`
	Books                []Book `xml:"books>book"`
}

// AuthorShow returns author information given a Goodreads author ID.
func (client Client) AuthorShow(ctx context.Context, id int) (Author, error) {
	type goodreadsResponse struct {
		Author Author `xml:"author"`
	}

	url := fmt.Sprintf("%s/author/show/%d.xml", client.getURL(), id)

	request, err := client.newRequestWithKey(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Author{}, err
	}

	response, err := client.getClient().Do(request)
	if err != nil {
		return Author{}, fmt.Errorf("do request: %w", err)
	}

	defer closeIgnoreError(response.Body)

	if response.StatusCode != http.StatusOK {
		return Author{}, ErrUnexpectedResponse{Code: response.StatusCode}
	}

	var author goodreadsResponse
	if err := xml.NewDecoder(response.Body).Decode(&author); err != nil {
		return Author{}, fmt.Errorf("decode response: %w", err)
	}

	return author.Author, nil
}
