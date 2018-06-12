package crawler

import (
	"webcrawler/workerpool"
)

type crawlerJob struct {
	url string
}

type crawlerJobResult struct {
	urls []string
	job  *crawlerJob
}

func (job *crawlerJob) Process() workerpool.JobResult {
	obtainedURLs := pageFetcher.Fetch(job.url)
	result := &crawlerJobResult{urls: obtainedURLs, job: job}
	return result
}

func (result *crawlerJobResult) GetJob() workerpool.Job {
	return result.job
}
