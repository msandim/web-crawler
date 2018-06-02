package crawler

import (
	"webcrawler/workerpool"
)

// Crawler is
type Crawler struct {
	// Parameters regarding the pool, with access to the jobs and results channel
	pool        *workerpool.WorkerPool
	pendingURLs chan workerpool.Job
	results     chan workerpool.JobResult

	// Parameters related to the crawling process:
	domain      string
	maxDepth    int
	checkedUrls map[string]bool
}

// New generates a Crawler struct:
func New(nWorkers int, domain string, maxDepth int) *Crawler {
	jobs := make(chan workerpool.Job, 10)
	results := make(chan workerpool.JobResult, 10)
	return &Crawler{
		pool:        workerpool.New(nWorkers, jobs, results),
		pendingURLs: jobs,
		results:     results,
		domain:      domain,
		maxDepth:    maxDepth,
		checkedUrls: make(map[string]bool),
	}
}

// Run initiates the crawler by running its routine "onJobProcessed" and the Worker Pool.
// It also adds the first crawling task: the domain's page.
func (crawler *Crawler) Run() {
	crawler.pool.Run()

	// Start the first job: crawl the main page of the domain:
	crawler.pool.AddJob(&crawlerJob{
		url:  crawler.domain,
		pool: crawler.pool,
	})
	crawler.checkedUrls[crawler.domain] = true

	go onURLCrawled(crawler)
}

// onUrlCrawled is a routine that iterates over the results returned by the Worker Pool
// and generates new crawling tasks for the Workers.
// In this case, new urls to crawl that haven't been checked before
func onURLCrawled(crawler *Crawler) {
	for result := range crawler.results {
		jobResult := result.(*crawlerJobResult)

		// Iterate over the URLs on the page we obtained:
		for _, url := range jobResult.urls {

			// If we never crawled that url, then we do it now:
			if !crawler.checkedUrls[url] {
				crawler.pool.AddJob(&crawlerJob{url: url})
				crawler.checkedUrls[url] = true
			}
		}
	}
}
