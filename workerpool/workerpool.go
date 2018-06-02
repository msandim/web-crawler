package workerpool

import (
	"fmt"
	"sync"
)

// Job is
type Job interface {
	Process() JobResult
}

// JobResult is
type JobResult interface {
	Process()
}

// WorkerPool is
type WorkerPool struct {
	nWorkers   int
	jobs       chan Job
	jobResults chan JobResult
	jobsActive *sync.WaitGroup
}

var mutex sync.Mutex

// New generates a WorkerPool struct and runs "nWorkers" workers
func New(nWorkers int, jobs chan Job, jobResults chan JobResult) *WorkerPool {

	pool := &WorkerPool{
		nWorkers:   nWorkers,
		jobs:       jobs,
		jobResults: jobResults,
		jobsActive: &sync.WaitGroup{},
	}

	return pool
}

// Run initiates the Worker Pool:
func (pool *WorkerPool) Run() {
	// Create a go routine for each worker:
	for i := 0; i < pool.nWorkers; i++ {
		go workerRoutine(pool)
	}

	// Create a go routine to process results:
	go resultRoutine(pool)
}

func workerRoutine(pool *WorkerPool) {
	// While there are jobs to process:
	for job := range pool.jobs {
		result := job.Process()
		pool.jobsActive.Done()
		pool.jobResults <- result
	}
}

func resultRoutine(pool *WorkerPool) {
	for result := range pool.jobResults {
		fmt.Println("Processed result from job", result)
	}

	//done <- true
}

// Wait waits until all the workers finished and returns:
/*
func (pool *WorkerPool) Wait() {
	pool.workersActive.Wait()
	close(pool.jobResults)
}
*/

// AddJob adds a job to the pool of workers:
func (pool *WorkerPool) AddJob(job Job) {
	mutex.Lock()
	pool.jobsActive.Add(1)

	select {
	case pool.jobs <- job: // some other worker can do it:
	default: // this routine will need to do that job:
		result := job.Process()
		pool.jobsActive.Done()
		pool.jobResults <- result
	}

	mutex.Unlock()
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
