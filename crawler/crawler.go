package crawler

import "webcrawler/workerpool"

type Crawler struct {
	// Parameters regarding the pool, with access to the jobs and results channel
	pool    *workerpool.WorkerPool
	jobs    chan workerpool.Job
	results chan workerpool.JobResult

	// Parameters related to the crawling process:
	domain   string
	maxDepth int
}

// New generates a Crawler struct:
func New(nWorkers int, domainArg string, maxDepthArg int) *Crawler {
	jobs := make(chan workerpool.Job, 10)
	results := make(chan workerpool.JobResult, 10)
	pool := workerpool.New(nWorkers, jobs, results)
	domain := domainArg
	maxDepth := maxDepthArg
	return &Crawler{pool: pool, jobs: jobs, results: results, domain: domain, maxDepth: maxDepth}
}

// Run initiates the crawler
func Run() {

}
