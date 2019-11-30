package goodreads_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/BooleanCat/go-goodreads"
	"github.com/BooleanCat/go-goodreads/fakes"
	"github.com/BooleanCat/go-goodreads/httpclient"
)

func ExampleClient_UserShow() {
	client := goodreads.Client{
		Client: httpclient.RateLimited(http.DefaultClient, ticker),
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
	if err != nil {
		t.Fatalf(`expected error "%v" not to have occurred`, err)
	}

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
	if user != want {
		t.Fatalf(`expected user "%v" to equal "%v"`, user, want)
	}

	if fakeDoer.DoCallCount() != 1 {
		t.Fatal("expected request to have been performed")
	}

	request := fakeDoer.DoArgsForCall(0)

	if request.Method != http.MethodGet {
		t.Fatalf(`expected request method to equal "%s"`, http.MethodGet)
	}

	wantURL := "https://www.goodreads.com/user/show/foo.xml?key=key"
	if request.URL.String() != wantURL {
		t.Fatalf(`expected URL "%s" to equal %s`, request.URL.String(), wantURL)
	}
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
	if err == nil {
		t.Fatal("expected failure")
	}

	if !strings.Contains(err.Error(), "create request") {
		t.Fatalf(`expected error "%s" to contain "create request"`, err.Error())
	}

	if fakeDoer.DoCallCount() != 0 {
		t.Fatal("expected request not to have been performed")
	}
}

func TestClient_UserShow_DoRequestFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(nil, errors.New("oops"))
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.UserShow("foo")
	if err == nil {
		t.Fatal("expected failure")
	}

	if err.Error() != "do request: oops" {
		t.Fatalf(`expected error "%s" to equal "do request: oops"`, err.Error())
	}
}

func TestClient_UserShow_InvalidStatusCode(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusMethodNotAllowed,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.UserShow("foo")
	if err == nil {
		t.Fatal("expected failure")
	}

	want := `unexpected status code "405"`
	if err.Error() != want {
		t.Fatalf(`expected error "%s" to equal "%s`, err.Error(), want)
	}
}

func TestClient_UserShow_DecodeFails(t *testing.T) {
	fakeDoer := new(fakes.FakeDoer)
	fakeDoer.DoReturns(&http.Response{
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		StatusCode: http.StatusOK,
	}, nil)
	client := goodreads.Client{Client: fakeDoer, Key: "key"}

	_, err := client.UserShow("foo")
	if err == nil {
		t.Fatal("expected failure")
	}

	if !strings.Contains(err.Error(), "decode response") {
		t.Fatalf(`expected error "%s" to contain "decode response"`, err.Error())
	}
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
