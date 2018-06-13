package crawler

import (
	"webcrawler/fetcher/URLWrapper"
	"webcrawler/workerpool"
)

// Implementation of the Crawling Jobs for the Worker Pool:

type crawlerJob struct {
	url string
}

type crawlerJobResult struct {
	urls []string
	job  *crawlerJob
}

func (job *crawlerJob) Process() workerpool.JobResult {
	obtainedURLs, errs := pageFetcher.Fetch(urlwrapper.New(job.url))

	for _, err := range errs {
		log.logError(err.Error())
	}

	result := &crawlerJobResult{urls: obtainedURLs, job: job}
	return result
}

func (result *crawlerJobResult) GetJob() workerpool.Job {
	return result.job
}
