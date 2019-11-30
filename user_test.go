package goodreads_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/BooleanCat/go-goodreads"
	"github.com/BooleanCat/go-goodreads/assert"
	"github.com/BooleanCat/go-goodreads/fakes"
	"github.com/BooleanCat/go-goodreads/httputils"
)

func ExampleClient_UserShow() {
	client := goodreads.Client{
		Client: httputils.DripLimit(http.DefaultClient, ticker),
	}

	user, err := client.UserShow("101333864")
	if err != nil {
		panic(err)
	}

	fmt.Println(user.UserName)
	// Output:
	// tgodkin
}

func TestClient_UserShow(t *testing.T) {
	responseBody := bytes.NewBufferString(userShowResponseBody)

	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(responseBody),
		StatusCode: http.StatusOK,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	user, err := client.UserShow("foo")
	assert.Nil(t, err)

	want := goodreads.User{
		ID:            "foo",
		Name:          "Foo Bar",
		UserName:      "fbar",
		Link:          "https://foo.com/fbar",
		ImageURL:      "https://foo.com/fbar.png",
		SmallImageURL: "https://foo.com/fbarmini.png",
		About:         "A test user",
		Age:           "30",
		Gender:        "male",
		Location:      "London, The United Kingdom",
		Website:       "https://foo.com",
	}
	assert.Equal(t, user, want)

	assert.Equal(t, fakeDoer.DoCallCount(), 1)
	request := fakeDoer.DoArgsForCall(0)
	assert.Equal(t, request.Method, http.MethodGet)
	assert.Equal(t, request.URL.String(), "https://www.goodreads.com/user/show/foo.xml?key=key")
}

func TestClient_UserShow_CreateRequestFails(t *testing.T) {
	responseBody := bytes.NewBufferString(userShowResponseBody)

	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(responseBody),
		StatusCode: http.StatusOK,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.UserShow("%%%%%%")
	assert.ErrorMatches(t, err, `^create request: `)
	assert.Equal(t, fakeDoer.DoCallCount(), 0)
}

func TestClient_UserShow_DoRequestFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(nil, errors.New("oops"))
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.UserShow("foo")
	assert.ErrorMatches(t, err, `^do request: oops$`)
}

func TestClient_UserShow_InvalidStatusCode(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusMethodNotAllowed,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.UserShow("foo")
	assert.ErrorMatches(t, err, `^unexpected status code "405"$`)
}

func TestClient_UserShow_DecodeFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusOK,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.UserShow("foo")
	assert.ErrorMatches(t, err, `^decode response: `)
}

const userShowResponseBody string = `
	<goodreads_response>
		<user>
			<id>foo</id>
			<name>Foo Bar</name>
			<user_name>fbar</user_name>
			<link><![CDATA[https://foo.com/fbar]]></link>
			<image_url><![CDATA[https://foo.com/fbar.png]]></image_url>
			<small_image_url><![CDATA[https://foo.com/fbarmini.png]]></small_image_url>
			<about>A test user</about>
			<age>30</age>
			<gender>male</gender>
			<location>London, The United Kingdom</location>
			<website>https://foo.com</website>
		</user>
	</goodreads_response>
`
