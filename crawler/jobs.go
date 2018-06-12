package crawler

import (
	"fmt"
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
	obtainedURLs, err := pageFetcher.Fetch(job.url)

	if err != nil {
		fmt.Println(err)
	}

	result := &crawlerJobResult{urls: obtainedURLs, job: job}
	return result
}

func (result *crawlerJobResult) GetJob() workerpool.Job {
	return result.job
}
