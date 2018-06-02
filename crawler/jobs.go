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
}

func (job *crawlerJobResult) Process() {
	fmt.Println("Resultado processado")
	//time.Sleep(2 * time.Second)
}

func (job *crawlerJob) Process() workerpool.JobResult {
	fmt.Println("Job processado")
	//time.Sleep(3 * time.Second)
	return &crawlerJobResult{}
}
