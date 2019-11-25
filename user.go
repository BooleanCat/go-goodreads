package goodreads

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type User struct {
	ID       string `xml:"id"`
	Name     string `xml:"name"`
	UserName string `xml:"user_name"`
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

	response, err := client.client.Do(request)
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
