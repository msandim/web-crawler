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
	GetJob() Job
	Process()
}

// WorkerPool is
type WorkerPool struct {
	nWorkers     int
	pendingJobs  chan Job
	finishedJobs chan JobResult
	jobsActive   *sync.WaitGroup
}

var nJobsActive sync.Mutex

// New generates a WorkerPool struct and runs "nWorkers" workers
func New(nWorkers int, jobs chan Job, jobResults chan JobResult) *WorkerPool {

	pool := &WorkerPool{
		nWorkers:     nWorkers,
		pendingJobs:  jobs,
		finishedJobs: jobResults,
		jobsActive:   &sync.WaitGroup{},
	}

	return pool
}

// Run initiates the Worker Pool:
func (pool *WorkerPool) Run() {
	// Create a go routine for each worker:
	for i := 0; i < pool.nWorkers; i++ {
		go workerRoutine(pool)
	}

	go checkEndRoutine(pool)
}

// AddJob adds a job to the pool of workers:
func (pool *WorkerPool) AddJob(job Job) {

	pool.jobsActive.Add(1)

	select {
	case pool.pendingJobs <- job: // some other worker can do it:

	default: // if the channel is full, do the job synchronously
		result := job.Process()
		result.Process()
		pool.jobsActive.Done()
	}
}

func workerRoutine(pool *WorkerPool) {
	// While there are jobs to process:
	for job := range pool.pendingJobs {
		// Add more jobs to chan jobs:
		result := job.Process()
		result.Process()
		pool.finishedJobs <- result
		pool.jobsActive.Done()
	}
}

func checkEndRoutine(pool *WorkerPool) {
	pool.jobsActive.Wait()
	fmt.Println("End of work, closing")
	close(pool.pendingJobs)
}
