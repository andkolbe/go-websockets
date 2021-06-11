package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// table test

// holds whatever we are posting to a page
type postData struct {
	// key   string
	// value string
}

var theTests = []struct { // slice of structs because we will have more than one test we want to run in our table test
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"login-page", "/", "GET", []postData{}, http.StatusOK},
	{"chat-page", "/chat", "GET", []postData{}, http.StatusOK},
	{"register-page", "/register", "GET", []postData{}, http.StatusOK},
	{"non-existent", "/nonexistent", "GET", []postData{}, http.StatusNotFound},
	{"login", "/", "POST", []postData{}, http.StatusNotFound},
	{"logout", "/logout", "POST", []postData{}, http.StatusNotFound},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes() // getRoutes comes from setup_test

	// we need to create a server and a client that can call that server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		// determine whether we are going to run a GET or POST request
		if e.method == "GET" {
			// we want to test as if we are a client accessing a web page
			response, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			// look at the status code we get back from the server and compare that against the expected status code
			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, response.StatusCode)
			}
		}
	}
}

var loginTests = []struct {
	name string
	username string
	password string
	expectedStatusCode int
	expectedHTML string
	expectedLocation string
}{
	{
		"valid-credentials",
		"test",
		"",
		http.StatusSeeOther,
		"",
		"/chat",
	},
	{
		"invalid-credentials",
		"invalidUsername",
		"",
		http.StatusSeeOther,
		"",
		"/",
	},
	// write an invalid data test to check for password validation when I implement it
	// {
	// 	"invalid-data",
	// 	"",
	// 	"p",
	// 	http.StatusOK, // status OK because the page rendered correctly, it just rendered a different page than where the user wanted to go
	// 	`action="/"`,
	// 	"/",
	// },
}

func TestLogin(t *testing.T) {
	// loop through all the tests
	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("username", e.username)
		postedData.Add("password", "password")

		// create request
		req, _ := http.NewRequest("POST", "/", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.Login)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: exepcted code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
