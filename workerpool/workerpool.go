package workerpool

import (
	"fmt"
	"sync"
)

type job interface {
	process() jobResult
}

type jobResult interface {
	process()
}

// WorkerPool is blablabla
type WorkerPool struct {
	nWorkers   int
	jobs       chan job
	jobResults chan jobResult
	waitGroup  *sync.WaitGroup
}

// New generates a WorkerPool struct and runs "nWorkers" workers
func New(nWorkers int) *WorkerPool {

	pool := &WorkerPool{
		nWorkers:   nWorkers,
		jobs:       make(chan job, nWorkers),
		jobResults: make(chan jobResult, nWorkers),
		waitGroup:  &sync.WaitGroup{},
	}

	// Create a go routine for each worker:
	for i := 0; i < nWorkers; i++ {
		pool.waitGroup.Add(1)
		go workerRoutine(pool)
	}

	// Create a go routine to process results:
	go resultRoutine(pool)

	return pool
}

func workerRoutine(pool *WorkerPool) {

	// While there are jobs to process:
	for job := range pool.jobs {
		result := job.process()
		fmt.Println(result)
	}

	// Announce that this worker finished:
	pool.waitGroup.Done()
}

func resultRoutine(pool *WorkerPool) {
	for result := range pool.jobResults {
		fmt.Println("Processed result from job", result)
	}

	//done <- true
}

// Wait waits until all the workers finished and returns:
func (pool *WorkerPool) Wait() {
	pool.waitGroup.Wait()
	close(pool.jobResults)
}

/*
func (pool *WorkerPool) AddJob(job Job) {

}
*/

/*
func worker(wg *sync.WaitGroup) {
	for job := range jobs {
		output := Result{job, digits(job.randomno)}
		results <- output
	}
	wg.Done()
}*/
