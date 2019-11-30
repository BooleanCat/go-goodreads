package goodreads

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type Book struct {
	ID                 string `xml:"id"`
	Title              string `xml:"title"`
	ISBN               string `xml:"isbn"`
	ISBN13             string `xml:"isbn13"`
	ASIN               string `xml:"asin"`
	KindleASIN         string `xml:"kindle_asin"`
	MarketplaceID      string `xml:"marketplace_id"`
	CountryCode        string `xml:"country_code"`
	ImageURL           string `xml:"image_url"`
	SmallImageURL      string `xml:"small_image_url"`
	PublicationYear    string `xml:"publication_year"`
	PublicationMonth   string `xml:"publication_month"`
	PublicationDay     string `xml:"publication_day"`
	Publisher          string `xml:"publisher"`
	LanguageCode       string `xml:"language_code"`
	IsEbook            string `xml:"is_ebook"`
	Description        string `xml:"description"`
	AverageRating      string `xml:"average_rating"`
	NumPages           string `xml:"num_pages"`
	Format             string `xml:"format"`
	EditionInformation string `xml:"edition_information"`
	RatingsCount       string `xml:"ratings_count"`
	TextReviewsCount   string `xml:"text_reviews_count"`
	URL                string `xml:"url"`
	Link               string `xml:"link"`
}

type Books struct {
	Books []Book `xml:"book"`
}

func (client Client) BookShow(id string) (Book, error) {
	type goodreadsResponse struct {
		Book Book `xml:"book"`
	}

	url := fmt.Sprintf("%s/book/show/%s.xml", goodreadsURL, id)
	response, err := client.doNewRequestWithKey(http.MethodGet, url, nil)
	if err != nil {
		return Book{}, err
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
