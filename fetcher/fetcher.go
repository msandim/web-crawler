package fetcher

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/html"
)

// Fetcher represents an entity that knows of to fetch the URLs
// contained in the HTML page of an URL
type Fetcher interface {
	Fetch(url string) []string
}

// HTTPFetcher implements the Fetcher interface and sends an HTTP GET to fetch
// the contents of an url
type HTTPFetcher struct {
	rateLimiter *RateLimiter
}

// NewHTTPFetcher returns a new HTTPFetcher with a given rate limit
// The rate limit corresponds to the number of concurrent requests
// that can be done
func NewHTTPFetcher(rateLimit int) *HTTPFetcher {
	return &HTTPFetcher{rateLimiter: NewRateLimiter(rateLimit)}
}

// Fetch sends an HTTP GET to fetch the contents of an url:
func (fetcher *HTTPFetcher) Fetch(urlArg string) []string {
	// URLs found in this page: avoid duplicates
	urlsFound := make(map[string]bool)

	// Parse the url we're trying to crawl, by extracting its url and path without url fragments:
	originalURLParsed, err := url.Parse(urlArg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "HTTPFetcher::fetch() - Error: failed to parse the URL to fetch: ", urlArg)
		return []string{}
	}

	// Define a custom http client that has a timeout and get the HTML code:
	var httpClient = &http.Client{Timeout: 10 * time.Second}
	fetcher.rateLimiter.Limit() // limit number of GET requests to be done at the same time
	resp, err := httpClient.Get(urlArg)
	fetcher.rateLimiter.Free()
	if err != nil {
		fmt.Fprintln(os.Stderr, "HTTPFetcher::fetch() - Error: Failed to GET: ", urlArg)
		return []string{}
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "HTTPFetcher::fetch() - Error: Failed to GET: ", urlArg, " with error code: ", resp.Status)
		return []string{}
	}

	body := resp.Body
	defer body.Close() // Close body when finishing reading from it
	tokenizer := html.NewTokenizer(body)

	for {
		tokenType := tokenizer.Next()

		switch {
		case tokenType == html.ErrorToken: // Reached the end of the document
			//fmt.Fprintln(os.Stderr, "HTTPFetcher::fetch() - Reached the end of page: ", urlArg)
			return mapToSlice(urlsFound)
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()

			// Check if the token corresponds to a <a>:
			if token.Data != "a" {
				continue
			}

			// Extract the href value, if there is one:
			extractedURL, ok := getHref(token)
			if !ok {
				fmt.Fprintln(os.Stderr, "HTTPFetcher::fetch() - Warning: <a> detected but no href present")
				continue
			}

			extractedURLParsed, err := url.Parse(extractedURL)
			if err != nil {
				fmt.Fprintln(os.Stderr, "HTTPFetcher::fetch() - Error: failed to parse the URL found: ", urlArg)
				continue
			}

			// Only crawl this new URL its domain is empty ("/otherpage") or if the domain of the url is the same:
			if extractedURLParsed.Hostname() == "" || extractedURLParsed.Hostname() == originalURLParsed.Hostname() {
				extractedURLParsed.Host = originalURLParsed.Host
				extractedURLParsed.Fragment = "" // delete fragments (e.g. #paragraph1)
				extractedURLParsed.RawQuery = "" // delete queries (?lang=en)

				if extractedURLParsed.Scheme == "" {
					extractedURLParsed.Scheme = "http"
				}

				// Only add to the map of found urls if we didn't add before:
				if _, ok := urlsFound[extractedURLParsed.String()]; !ok {
					urlsFound[extractedURLParsed.String()] = true
				}
			}

		}
	}
}

// getHref gets the href attribute from an <a> token
func getHref(token html.Token) (url string, ok bool) {
	// Iterate over all of the Token's attributes until we find an "href":
	for _, v := range token.Attr {
		if v.Key == "href" {
			return v.Val, true
		}
	}
	return "", false
}

func mapToSlice(urlMap map[string]bool) (urls []string) {
	urls = make([]string, 0, len(urlMap))
	for u := range urlMap {
		urls = append(urls, u)
	}
	return
}
