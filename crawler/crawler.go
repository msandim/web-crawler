package crawler

import (
	"web-crawler/fetcher"
	"web-crawler/workerpool"
)

// Crawler is a strucutre that contains the specifications of the crawling process.
type Crawler struct {
	// Parameters regarding the pool, with access to the jobs and results channel
	pool    *workerpool.WorkerPool
	results chan workerpool.JobResult

	// Parameters related to the crawling process:
	domain string

	// Variables for the crawler's state:
	nURLsCrawled int             // number of URLs successfully crawled
	checkedUrls  map[string]bool // number of URLs in which we initiated the crawling process
	finishedFlag chan bool
}

// Responsable to know how to fetch a page (through HTTP requests in production or mocked in testing):
var pageFetcher fetcher.Fetcher

// Responsable to know how to log the results:
var log logger = &printer{}

// New creates a Crawler struct given the arguments and returns a pointer to it.
func New(nWorkers int, rateLimit int, timeoutSeconds int, domain string) *Crawler {
	pageFetcher = fetcher.NewHTTPFetcher(rateLimit, timeoutSeconds)
	pool := workerpool.New(nWorkers)

	return &Crawler{
		pool:         pool,
		results:      pool.GetResultsChannel(),
		domain:       domain,
		checkedUrls:  make(map[string]bool),
		finishedFlag: make(chan bool),
	}
}

// newTesting creates Crawler struct given the arguments and returns a pointer to it (used only for testing).
func newTesting(nWorkers int, domain string) *Crawler {
	pool := workerpool.New(nWorkers)

	return &Crawler{
		pool:         pool,
		results:      pool.GetResultsChannel(),
		domain:       domain,
		checkedUrls:  make(map[string]bool),
		finishedFlag: make(chan bool),
	}
}

// Run initiates the crawler by running its routine "onJobProcessed" and the Worker Pool.
// It also adds the first crawling task: the domain's page.
// This function returns when the crawling process ended
func (crawler *Crawler) Run() {
	crawler.pool.Run()

	// Start the first job: crawl the main page of the domain:
	crawler.pool.AddJob(&crawlerJob{url: crawler.domain})
	crawler.checkedUrls[crawler.domain] = true

	// Initiate routine that will receive the crawling results:
	go onURLCrawled(crawler)

	// Wait for end of crawling process:
	<-crawler.finishedFlag
}

// onUrlCrawled is a routine that iterates over the results returned by the Worker Pool
// and generates new crawling tasks for the Workers.
// In this case, new urls to crawl that haven't been checked before.
func onURLCrawled(crawler *Crawler) {
	for result := range crawler.results {

		job := result.GetJob().(*crawlerJob)
		parentURL := job.url
		childrenURLs := []string{}

		// Get the result from crawling job and increment the number of URLs crawled:
		jobResult := result.(*crawlerJobResult)
		crawler.nURLsCrawled++

		// Iterate over the URLs on the page we obtained:
		for _, url := range jobResult.urls {

			childrenURLs = append(childrenURLs, url)

			// If we never crawled that url, then we do it now:
			if !crawler.checkedUrls[url] {
				crawler.pool.AddJob(&crawlerJob{url: url})
				crawler.checkedUrls[url] = true
			}
		}

		log.logPage(parentURL, childrenURLs)

		// if all the URLs launched for crawling had their crawling processes ended:
		if len(crawler.checkedUrls) == crawler.nURLsCrawled {
			crawler.pool.EndJobs()
		}
	}

	crawler.finishedFlag <- true
}
