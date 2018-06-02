package crawler

import (
	"fmt"
	"webcrawler/workerpool"
)

type crawlerJob struct {
	url  string
	pool *workerpool.WorkerPool
}

type crawlerJobResult struct {
	urls []string
	job  *crawlerJob
}

func (job *crawlerJob) Process() workerpool.JobResult {
	fmt.Println("Job processado")
	result := crawlerJobResult{job: job}
	return &result
}

func (result *crawlerJobResult) Process() {
	fmt.Println("Resultado processado")
}

func (result *crawlerJobResult) GetJob() workerpool.Job {
	return result.job
}
