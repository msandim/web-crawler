package fetcher

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/html"
)

// Fetcher represents an entity that knows of to fetch the HTML contents of
// an url
type Fetcher interface {
	Fetch(url string) string
}

// HTTPFetcher implements the Fetcher interface and sends an HTTP GET to fetch
// the contents of an url
type HTTPFetcher struct {
}

// Fetch sends an HTTP GET to fetch the contents of an url:
func (fetcher *HTTPFetcher) Fetch(url string) []string {
	var urls []string

	// Get the domain of the url we need to crawl:
	domain, err := getDomain(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return urls
	}

	// Define a custom http client that has a timeout and get the HTML code:
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "HTTPFetcher::fetch() - Error: Failed to GET: \""+url+"\"")
		return urls
	}

	body := resp.Body
	defer body.Close() // Close body when finishing reading from it
	tokenizer := html.NewTokenizer(body)

	for {
		tokenType := tokenizer.Next()

		switch {
		case tokenType == html.ErrorToken: // Reached the end of the document
			return urls
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()

			// Check if the token corresponds to a <a>:
			if token.Data != "a" {
				continue
			}

			// Extract the href value, if there is one:
			url, ok := getHref(token)
			if !ok {
				fmt.Fprintln(os.Stderr, "HTTPFetcher::fetch() - Warning: <a> detected but no href present")
				continue
			}

			// Check if the url belongs to the domain:
			if !belongsToDomain(url, domain) {
				continue
			}

			urls = append(urls, url)
		}
	}
}

// Helper function that gets the href attribute from an <a> token
func getHref(token html.Token) (url string, ok bool) {
	// Iterate over all of the Token's attributes until we find an "href":
	for _, v := range token.Attr {
		if v.Key == "href" {
			return v.Val, true
		}
	}
	return "", false
}

func getDomain(urlArg string) (string, error) {
	u, err := url.Parse(urlArg)
	if err != nil {
		return "", errors.New("getDomain() - Error: failed to get domain from url (" + urlArg + ")")
	}
	return u.Hostname(), nil
}

// belongsToDomain checks if the domain of the "url" argument is the same as "domain"
// if the "url" argument has an invalid domain, it returns false
func belongsToDomain(url, domain string) bool {
	urlDomain, err := getDomain(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	return domain == urlDomain
}
