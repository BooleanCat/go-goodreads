package goodreads

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type Author struct {
	ID                   string `xml:"id"`
	Name                 string `xml:"name"`
	Link                 string `xml:"link"`
	FansCount            int64  `xml:"fans_count"`
	AuthorFollowersCount int64  `xml:"author_followers_count"`
	LargeImageURL        string `xml:"large_image_url"`
	ImageURL             string `xml:"image_url"`
	SmallImageURL        string `xml:"small_image_url"`
	About                string `xml:"about"`
	Influences           string `xml:"influences"`
	WorksCount           string `xml:"works_count"`
	Gender               string `xml:"gender"`
	Hometown             string `xml:"hometown"`
	BornAt               string `xml:"born_at"`
	DiedAt               string `xml:"died_at"`
	GoodreadsAuthor      string `xml:"goodreads_author"`
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
