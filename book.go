package goodreads

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type Book struct {
	ID                 int          `xml:"id"`
	Title              string       `xml:"title"`
	ISBN               string       `xml:"isbn"`
	ISBN13             string       `xml:"isbn13"`
	ASIN               string       `xml:"asin"`
	KindleASIN         string       `xml:"kindle_asin"`
	MarketplaceID      string       `xml:"marketplace_id"`
	CountryCode        string       `xml:"country_code"`
	ImageURL           string       `xml:"image_url"`
	SmallImageURL      string       `xml:"small_image_url"`
	PublicationYear    int          `xml:"publication_year"`
	PublicationMonth   int          `xml:"publication_month"`
	PublicationDay     int          `xml:"publication_day"`
	Publisher          string       `xml:"publisher"`
	LanguageCode       string       `xml:"language_code"`
	IsEbook            bool         `xml:"is_ebook"`
	Description        string       `xml:"description"`
	AverageRating      float32      `xml:"average_rating"`
	NumPages           int          `xml:"num_pages"`
	Format             string       `xml:"format"`
	EditionInformation string       `xml:"edition_information"`
	RatingsCount       int          `xml:"ratings_count"`
	TextReviewsCount   int          `xml:"text_reviews_count"`
	URL                string       `xml:"url"`
	Link               string       `xml:"link"`
	Work               Work         `xml:"work"`
	Authors            []Author     `xml:"authors>author"`
	PopularShelves     []Shelf      `xml:"popular_shelves>shelf"`
	BookLinks          []Link       `xml:"book_links>book_link"`
	BuyLinks           []Link       `xml:"buy_links>buy_link"`
	SeriesWorks        []SeriesWork `xml:"series_works>series_work"`
	SimilarBooks       []Book       `xml:"similar_books>book"`
}

func (client Client) BookShow(id int) (Book, error) {
	type goodreadsResponse struct {
		Book Book `xml:"book"`
	}

	url := fmt.Sprintf("%s/book/show/%d.xml", goodreadsURL, id)
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

type Work struct {
	ID                             int64  `xml:"id"`
	BooksCount                     int64  `xml:"books_count"`
	BestBookID                     int64  `xml:"best_book_id"`
	ReviewsCount                   int64  `xml:"reviews_count"`
	RatingsSum                     int64  `xml:"ratings_sum"`
	RatingsCount                   int64  `xml:"ratings_count"`
	TextReviewsCount               int64  `xml:"text_reviews_count"`
	OriginalPublicationYear        int64  `xml:"original_publication_year"`
	OriginalPublicationMonth       int64  `xml:"original_publication_month"`
	OriginalPublicationDay         int64  `xml:"original_publication_day"`
	OriginalTitle                  string `xml:"original_title"`
	OriginalLanguageID             int64  `xml:"original_language_id"`
	MediaType                      string `xml:"media_type"`
	RatingDist                     string `xml:"rating_dist"`
	DescUserID                     int64  `xml:"desc_user_id"`
	DefaultChapteringBookID        int64  `xml:"default_chaptering_book_id"`
	DefaultDescriptionLanguageCode string `xml:"default_description_language_code"`
	WorkURI                        string `xml:"work_uri"`
}

type Shelf struct {
	Name  string `xml:"name,attr"`
	Count string `xml:"count,attr"`
}

type Link struct {
	ID   int    `xml:"id"`
	Name string `xml:"name"`
	Link string `xml:"link"`
}

type SeriesWork struct {
	ID           int    `xml:"id"`
	UserPosition int    `xml:"user_position"`
	Series       Series `xml:"series"`
}

type Series struct {
	ID               int    `xml:"id"`
	Title            string `xml:"title"`
	Description      string `xml:"description"`
	Note             string `xml:"note"`
	SeriesWorksCount int    `xml:"series_works_count"`
	PrimaryWorkCount int    `xml:"primary_work_count"`
	Numbered         bool   `xml:"numbered"`
}
