package fetcher

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"web-crawler/fetcher/URLWrapper"
)

func TestHTTPFetcher_Fetch_InvalidURL(t *testing.T) {
	domain := "%"
	errorMsg := "HTTPFetcher::fetch() - Error: failed to parse the URL to fetch: " + domain

	fetcher := NewHTTPFetcher(4, 10)
	urls, errs := fetcher.Fetch(urlwrapper.New(domain))

	if len(urls) != 0 {
		t.Errorf("Length of URLs was invalid. Expected: %d, Got: %d", 0, len(urls))
	}

	if len(errs) != 1 {
		t.Errorf("Length of errors was invalid. Expected: %d, Got: %d", 1, len(errs))
	}

	if errs[0].Error() != errorMsg {
		t.Errorf("Error message was not valid\nExpected: %s, Got: %s", errorMsg, errs[0])
	}
}

func TestHTTPFetcher_Fetch_InvalidGET(t *testing.T) {
	domain := "123"
	errorMsg := "HTTPFetcher::fetch() - Error: Failed to GET: " + domain

	fetcher := NewHTTPFetcher(4, 10)
	urls, errs := fetcher.Fetch(urlwrapper.New(domain))

	if len(urls) != 0 {
		t.Errorf("Length of URLs was invalid. Expected: %d, Got: %d", 0, len(urls))
	}

	if len(errs) != 1 {
		t.Errorf("Length of errors was invalid. Expected: %d, Got: %d", 1, len(errs))
	}

	if errs[0].Error() != errorMsg {
		t.Errorf("Error message was not valid\nExpected: %s, Got: %s", errorMsg, errs[0])
	}
}

func TestHTTPFetcher_Fetch_TimeoutGET(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer server.Close()

	domain := "http://monzo.com/"
	errorMsg := "HTTPFetcher::fetch() - Error: Failed to GET: " + domain

	fetcher := NewHTTPFetcher(4, 1)
	urls, errs := fetcher.Fetch(urlwrapper.NewTesting(domain, server.URL))

	if len(urls) != 0 {
		t.Errorf("Length of URLs was invalid. Expected: %d, Got: %d", 0, len(urls))
	}

	if len(errs) != 1 {
		t.Errorf("Length of errors was invalid. Expected: %d, Got: %d", 1, len(errs))
	}

	if errs[0].Error() != errorMsg {
		t.Errorf("Error message was not valid\nExpected: %s, Got: %s", errorMsg, errs[0])
	}
}

func TestHTTPFetcher_Fetch_ErrorMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	domain := "http://monzo.com/"
	errorMsg := "HTTPFetcher::fetch() - Error: Failed to GET: " + domain + " with error code: 404 Not Found"

	fetcher := NewHTTPFetcher(4, 10)
	urls, errs := fetcher.Fetch(urlwrapper.NewTesting(domain, server.URL))

	if len(urls) != 0 {
		t.Errorf("Length of URLs was invalid. Expected: %d, Got: %d", 0, len(urls))
	}

	if len(errs) != 1 {
		t.Errorf("Length of errors was invalid. Expected: %d, Got: %d", 1, len(errs))
	}

	if errs[0].Error() != errorMsg {
		t.Errorf("Error message was not valid\nExpected: %s, Got: %s", errorMsg, errs[0])
	}
}

func TestHTTPFetcher_Fetch_ContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/pdf")
	}))
	defer server.Close()

	domain := "http://monzo.com/"
	errorMsg := "HTTPFetcher::fetch() - Error: Content type of " + domain + " is application/pdf"

	fetcher := NewHTTPFetcher(4, 10)
	urls, errs := fetcher.Fetch(urlwrapper.NewTesting(domain, server.URL))

	if len(urls) != 0 {
		t.Errorf("Length of URLs was invalid. Expected: %d, Got: %d", 0, len(urls))
	}

	if len(errs) != 1 {
		t.Errorf("Length of errors was invalid. Expected: %d, Got: %d", 1, len(errs))
	}

	if errs[0].Error() != errorMsg {
		t.Errorf("Error message was not valid\nExpected: %s, Got: %s", errorMsg, errs[0])
	}
}

func TestHTTPFetcher_Fetch_Errors(t *testing.T) {

	page, filerr := ioutil.ReadFile("../test/example.html")
	if filerr != nil {
		t.Fatalf("HTML file not found")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))
	defer server.Close()

	domain := "http://monzo.com/"
	childURL := ""
	errorMsg := ""

	fetcher := NewHTTPFetcher(4, 10)
	urls, errs := fetcher.Fetch(urlwrapper.NewTesting(domain, server.URL))

	if len(urls) != 6 {
		t.Errorf("Length of URLs was invalid. Expected: %d, Got: %d", 6, len(urls))
	}

	childURL = "http://monzo.com/test1"
	if urls[0] != childURL {
		t.Errorf("Invalid child URL. Expected: %s, Got: %s", childURL, urls[0])
	}

	childURL = "http://monzo.com/test2"
	if urls[1] != childURL {
		t.Errorf("Invalid child URL. Expected: %s, Got: %s", childURL, urls[1])
	}

	childURL = "http://monzo.com/test3"
	if urls[2] != childURL {
		t.Errorf("Invalid child URL. Expected: %s, Got: %s", childURL, urls[2])
	}

	childURL = "http://monzo.com/test4"
	if urls[3] != childURL {
		t.Errorf("Invalid child URL. Expected: %s, Got: %s", childURL, urls[3])
	}

	childURL = "http://monzo.com/test5"
	if urls[4] != childURL {
		t.Errorf("Invalid child URL. Expected: %s, Got: %s", childURL, urls[4])
	}

	childURL = "https://monzo.com/test6"
	if urls[5] != childURL {
		t.Errorf("Invalid child URL. Expected: %s, Got: %s", childURL, urls[5])
	}

	if len(errs) != 2 {
		t.Errorf("Length of errors was invalid. Expected: %d, Got: %d", 2, len(errs))
	}

	errorMsg = "HTTPFetcher::fetch() - Warning: <a> detected but no href present"
	if errs[0].Error() != errorMsg {
		t.Errorf("Invalid error message. Expected: %s, Got: %s", childURL, errs[0])
	}

	errorMsg = "HTTPFetcher::fetch() - Warning: failed to parse the URL found: %"
	if errs[1].Error() != errorMsg {
		t.Errorf("Invalid error message. Expected: %s, Got: %s", childURL, errs[0])
	}

}
