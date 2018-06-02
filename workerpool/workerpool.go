package workerpool

import (
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

var nJobsActive sync.Mutex

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
}

// AddJob adds a job to the pool of workers:
func (pool *WorkerPool) AddJob(job Job) {
	pool.jobsActive.Add(1)

	select {
	case pool.jobs <- job: // some other worker can do it:
	default: // do the job synchronously
		result := job.Process()
		result.Process()
		pool.jobsActive.Done()
	}
}

func workerRoutine(pool *WorkerPool) {
	// While there are jobs to process:
	for job := range pool.jobs {
		// Add more jobs to chan jobs:
		result := job.Process()
		result.Process()
		pool.jobsActive.Done()
	}
}

/*
func resultRoutine(pool *WorkerPool) {
	for result := range pool.jobResults {
		result.Process()
	}
}
*/
