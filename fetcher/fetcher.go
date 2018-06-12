package fetcher

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Fetcher represents an entity that knows of to fetch the URLs
// contained in the HTML page of an URL
type Fetcher interface {
	Fetch(url string) ([]string, []error)
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
func (fetcher *HTTPFetcher) Fetch(urlArg string) ([]string, []error) {
	// URLs found in this page: avoid duplicates
	urlsFound := make(map[string]bool)
	errorsFound := []error{}

	// Parse the url we're trying to crawl, by extracting its url and path without url fragments:
	parentURLParsed, err := url.Parse(urlArg)
	if err != nil {
		errorsFound = append(errorsFound, errors.New("HTTPFetcher::fetch() - Error: failed to parse the URL to fetch: "+urlArg))
		return []string{}, errorsFound
	}

	// Define a custom http client that has a timeout and get the HTML code:
	var httpClient = &http.Client{Timeout: 10 * time.Second}
	fetcher.rateLimiter.Limit() // limit number of GET requests to be done at the same time
	resp, err := httpClient.Get(urlArg)
	fetcher.rateLimiter.Free()
	if err != nil {
		errorsFound = append(errorsFound, errors.New("HTTPFetcher::fetch() - Error: Failed to GET: "+urlArg))
		return []string{}, errorsFound
	}
	if resp.StatusCode != http.StatusOK {
		errorsFound = append(errorsFound, errors.New("HTTPFetcher::fetch() - Error: Failed to GET: "+urlArg+" with error code: "+resp.Status))
		return []string{}, errorsFound
	}

	defer resp.Body.Close() // Close body when finishing reading from it

	// Only proceed if it's an HTML document:
	if !strings.Contains(resp.Header.Get("Content-type"), "text/html") {
		errorsFound = append(errorsFound, errors.New("HTTPFetcher::fetch() - Error: Content type of "+urlArg+" is "+resp.Header.Get("Content-type")))
		return []string{}, errorsFound
	}

	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tokenType := tokenizer.Next()

		switch {
		case tokenType == html.ErrorToken: // Reached the end of the document
			return mapToSlice(urlsFound), nil
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()

			// Check if the token corresponds to a <a>:
			if token.Data != "a" {
				continue
			}

			// Extract the href value, if there is one:
			childURL, ok := getHref(token)
			if !ok {
				errorsFound = append(errorsFound, errors.New("HTTPFetcher::fetch() - Warning: <a> detected but no href present"))
				continue
			}

			childURLParsed, err := url.Parse(childURL)
			if err != nil {
				errorsFound = append(errorsFound, errors.New("HTTPFetcher::fetch() - Warning: failed to parse the URL found: "+urlArg))
				continue
			}

			if !isChildURLValid(childURLParsed, *parentURLParsed) {
				continue
			}

			// Only add to the map of found urls if we didn't add before:
			if _, ok := urlsFound[childURLParsed.String()]; !ok {
				urlsFound[childURLParsed.String()] = true
			}
		}
	}
}

func isChildURLValid(childURL *url.URL, fatherURL url.URL) bool {
	// Only crawl this new URL if its domain is empty (e.g. "/otherpage") or if the domain of the url is the same:
	if childURL.Hostname() == "" || childURL.Hostname() == fatherURL.Hostname() {
		childURL.Host = fatherURL.Host
		childURL.Fragment = "" // delete fragments (e.g. #paragraph1)
		childURL.RawQuery = "" // delete queries (?lang=en)

		// We're only interested in http, https and or in links with no scheme (avoid tel e mailto):
		if childURL.Scheme != "http" && childURL.Scheme != "https" && childURL.Scheme != "" {
			return false
		}

		if childURL.Scheme == "" {
			childURL.Scheme = "http"
		}

		return true
	}
	return false
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
