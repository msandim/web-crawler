package urlwrapper

// URLWrapper is a wrapper around a URL, that has a different URL for URL analysis
// and a different one to perform HTTP requests.
type URLWrapper struct {
	URL           string
	URLForRequest string
}

// New generates a new URLWrapper, in which the URL for HTTP requests is the same as the regular URL
func New(url string) *URLWrapper {
	return &URLWrapper{
		URL:           url,
		URLForRequest: url,
	}
}

// NewTesting generates a new URLWrapper, with a different URL for requests
func NewTesting(url string, urlForRequest string) *URLWrapper {
	return &URLWrapper{
		URL:           url,
		URLForRequest: urlForRequest,
	}
}
