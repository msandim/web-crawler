package crawler

import (
	"fmt"
	"webcrawler/workerpool"
)

// Crawler is
type Crawler struct {
	// Parameters regarding the pool, with access to the jobs and results channel
	pool    *workerpool.WorkerPool
	results chan workerpool.JobResult

	// Parameters related to the crawling process:
	domain   string
	maxDepth int

	// Variables for the crawler's state:
	nURLsCrawled int
	checkedUrls  map[string]bool
	finishedFlag chan bool
}

// New generates a Crawler struct:
func New(nWorkers int, domain string, maxDepth int) *Crawler {
	jobs := make(chan workerpool.Job, 10)
	results := make(chan workerpool.JobResult, 10)
	return &Crawler{
		pool: workerpool.New(nWorkers, jobs, results),
		//pendingURLs: jobs,
		results:      results,
		domain:       domain,
		maxDepth:     maxDepth,
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
// In this case, new urls to crawl that haven't been checked before
func onURLCrawled(crawler *Crawler) {
	fmt.Println("onURLCrawled() - vou começar a receber resultados")
	for result := range crawler.results {

		fmt.Println("onURLCrawled() - Processed: ", result.GetJob())

		// Get the result from crawling job and increment the number of URLs crawled:
		jobResult := result.(*crawlerJobResult)
		crawler.nURLsCrawled++

		// Iterate over the URLs on the page we obtained:
		for _, url := range jobResult.urls {

			// If we never crawled that url, then we do it now:
			if !crawler.checkedUrls[url] {
				fmt.Println("onURLCrawled() - Added for processing: ", url)
				crawler.pool.AddJob(&crawlerJob{url: url})
				crawler.checkedUrls[url] = true
			}
		}

		// if all the URLs launched for crawling had their crawling processes ended:
		if len(crawler.checkedUrls) == crawler.nURLsCrawled {
			crawler.pool.EndJobs()
			crawler.finishedFlag <- true
		}
	}
}
