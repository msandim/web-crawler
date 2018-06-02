package workerpool

import (
	"fmt"
	"sync"
)

// Job is
type Job interface {
	Process() JobResult
}

// JobGenerator is
/*
type JobGenerator interface {
	Generate(JobResult) []Job
}
*/

// JobResult is
type JobResult interface {
	Process()
}

// WorkerPool is
type WorkerPool struct {
	nWorkers   int
	jobs       chan Job
	jobResults chan JobResult
	waitGroup  *sync.WaitGroup
}

// New generates a WorkerPool struct and runs "nWorkers" workers
func New(nWorkers int, jobs chan Job, jobResults chan JobResult) *WorkerPool {

	pool := &WorkerPool{
		nWorkers:   nWorkers,
		jobs:       jobs,
		jobResults: jobResults,
		waitGroup:  &sync.WaitGroup{},
	}

	return pool
}

// Run initiates the Worker Pool:
func (pool *WorkerPool) lRun() {
	// Create a go routine for each worker:
	for i := 0; i < pool.nWorkers; i++ {
		pool.waitGroup.Add(1)
		go workerRoutine(pool)
	}

	// Create a go routine to process results:
	go resultRoutine(pool)
}

func workerRoutine(pool *WorkerPool) {
	// While there are jobs to process:

	for job := range pool.jobs {
		result := job.Process()
		fmt.Println("Processed job", result)

		pool.jobResults <- result
	}

	// Announce that this worker finished:
	fmt.Println("FIZ DONE")
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

// AddJob adds a job to the pool of workers:
func (pool *WorkerPool) AddJob(job Job) {
	pool.jobs <- job
}

//func (pool *WorkerPool) EndJobs

/*
func worker(wg *sync.WaitGroup) {
	for job := range jobs {
		output := Result{job, digits(job.randomno)}
		results <- output
	}
	wg.Done()
}*/
