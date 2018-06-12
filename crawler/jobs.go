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
	obtainedURLs, errs := pageFetcher.Fetch(job.url)

	for _, err := range errs {
		log.logError(err)
	}

	result := &crawlerJobResult{urls: obtainedURLs, job: job}
	return result
}

func (result *crawlerJobResult) GetJob() workerpool.Job {
	return result.job
}
