package goodreads

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type User struct {
	ID            string `xml:"id"`
	Name          string `xml:"name"`
	UserName      string `xml:"user_name"`
	Link          string `xml:"link"`
	ImageURL      string `xml:"image_url"`
	SmallImageURL string `xml:"small_image_url"`
	About         string `xml:"about"`
	Age           string `xml:"age"`
	Gender        string `xml:"gender"`
	Location      string `xml:"location"`
	Website       string `xml:"website"`
}

func (client Client) UserShow(id string) (User, error) {
	type goodreadsResponse struct {
		User User `xml:"user"`
	}

	url := fmt.Sprintf("%s/user/show/%s.xml", goodreadsURL, id)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return User{}, fmt.Errorf("create request: %w", err)
	}

	request, err = client.addGoodreadsKeyQueryParam(request)
	if err != nil {
		return User{}, err
	}

	response, err := client.Client.Do(request)
	if err != nil {
		return User{}, fmt.Errorf("do request: %w", err)
	}
	defer closeIgnoreError(response.Body)

	if response.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf(`unexpected status code "%d"`, response.StatusCode)
	}

	var user goodreadsResponse
	if err := xml.NewDecoder(response.Body).Decode(&user); err != nil {
		return User{}, fmt.Errorf("decode response: %w", err)
	}

	return user.User, nil
}
