package crawler

type Fetcher interface {
	fetch(url string) string
}
