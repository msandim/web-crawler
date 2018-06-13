package fetcher

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
	"webcrawler/fetcher/URLWrapper"

	"golang.org/x/net/html"
)

// Fetcher represents an entity that knows of to fetch the URLs
// contained in the HTML page of an URL.
type Fetcher interface {
	Fetch(urlArg *urlwrapper.URLWrapper) ([]string, []error)
}

// HTTPFetcher implements the Fetcher interface and sends an HTTP GET to fetch
// the contents of an url.
type HTTPFetcher struct {
	rateLimiter    *RateLimiter
	timeoutSeconds int
}

// NewHTTPFetcher returns a new HTTPFetcher with a given rate limit
// The rate limit corresponds to the number of concurrent requests
// that can be done.
func NewHTTPFetcher(rateLimit int, timeoutSeconds int) *HTTPFetcher {
	return &HTTPFetcher{
		rateLimiter:    NewRateLimiter(rateLimit),
		timeoutSeconds: timeoutSeconds,
	}
}

// Fetch sends an HTTP GET to fetch the contents of an url.
func (fetcher *HTTPFetcher) Fetch(urlArg *urlwrapper.URLWrapper) ([]string, []error) {
	// URLs found in this page: avoid duplicates:
	urlsFound := []string{}
	urlsFoundMap := make(map[string]bool)
	errorsFound := []error{}

	// Parse the url we're trying to crawl, by extracting its url and path without url fragments:
	parentURLParsed, err := url.Parse(urlArg.URL)
	if err != nil {
		errorsFound = append(errorsFound, errors.New("HTTPFetcher::fetch() - Error: failed to parse the URL to fetch: "+urlArg.URL))
		return []string{}, errorsFound
	}

	// Define a custom http client that has a timeout and get the HTML code:
	var httpClient = &http.Client{Timeout: time.Duration(fetcher.timeoutSeconds) * time.Second}

	fetcher.rateLimiter.Limit() // limit number of GET requests to be done at the same time
	resp, err := httpClient.Get(urlArg.URLForRequest)
	fetcher.rateLimiter.Free()

	if err != nil {
		errorsFound = append(errorsFound, errors.New("HTTPFetcher::fetch() - Error: Failed to GET: "+urlArg.URL))
		return []string{}, errorsFound
	}
	if resp.StatusCode != http.StatusOK {
		errorsFound = append(errorsFound, errors.New("HTTPFetcher::fetch() - Error: Failed to GET: "+urlArg.URL+" with error code: "+resp.Status))
		return []string{}, errorsFound
	}

	defer resp.Body.Close() // Close body when finishing reading from it

	// Only proceed if it's an HTML document:
	if !strings.Contains(resp.Header.Get("Content-type"), "text/html") {
		errorsFound = append(errorsFound, errors.New("HTTPFetcher::fetch() - Error: Content type of "+urlArg.URL+" is "+resp.Header.Get("Content-type")))
		return []string{}, errorsFound
	}

	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tokenType := tokenizer.Next()

		switch {
		case tokenType == html.ErrorToken: // Reached the end of the document
			return urlsFound, errorsFound
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
				errorsFound = append(errorsFound, errors.New("HTTPFetcher::fetch() - Warning: failed to parse the URL found: "+childURL))
				continue
			}

			if !isChildURLValid(childURLParsed, *parentURLParsed) {
				continue
			}

			// Only add to the map of found urls if we didn't add before:
			if _, ok := urlsFoundMap[childURLParsed.String()]; !ok {
				urlsFoundMap[childURLParsed.String()] = true
				urlsFound = append(urlsFound, childURLParsed.String())
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
