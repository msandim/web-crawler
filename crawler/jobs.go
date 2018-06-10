package crawler

import (
	"fmt"
	"webcrawler/fetcher"
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
	fmt.Println("crawlerJob::Process() - Starting to process: " + job.url)
	fetcher := &fetcher.HTTPFetcher{}
	obtainedURLs := fetcher.Fetch(job.url)
	result := &crawlerJobResult{urls: obtainedURLs, job: job}
	return result

	/* 	var result crawlerJobResult

	   	if job.url == "1" {
	   		result = crawlerJobResult{urls: []string{"2", "3"}}
	   	}

	   	if job.url == "2" {
	   		result = crawlerJobResult{urls: []string{"4", "5"}}
	   	}

	   	if job.url == "3" {
	   		result = crawlerJobResult{urls: []string{"6"}}
	   	}

	   	if job.url == "4" {
	   		result = crawlerJobResult{urls: []string{"2"}}
	   	}

	   	result.job = job

	   	fmt.Println("crawlerJob::Process() - Job processado")

	   	return &result */

}

func (result *crawlerJobResult) GetJob() workerpool.Job {
	return result.job
}

/*
func (result *crawlerJobResult) Process() {
	fmt.Println("Resultado processado")
}


*/
