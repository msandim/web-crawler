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
	GetJob() Job
}

// WorkerPool is
type WorkerPool struct {
	nWorkers      int
	pendingJobs   chan Job
	finishedJobs  chan JobResult
	workersActive *sync.WaitGroup
}

// New generates a WorkerPool struct and runs "nWorkers" workers.
func New(nWorkers int, jobs chan Job, jobResults chan JobResult) *WorkerPool {

	pool := &WorkerPool{
		nWorkers:      nWorkers,
		pendingJobs:   jobs,
		finishedJobs:  jobResults,
		workersActive: &sync.WaitGroup{},
	}

	return pool
}

// Run initiates the Worker Pool.
func (pool *WorkerPool) Run() {
	// Create a goroutine for each worker:
	for i := 0; i < pool.nWorkers; i++ {
		pool.workersActive.Add(1)
		go workerRoutine(pool)
	}

	go waitForWorkersRoutine(pool)
}

// AddJob adds a job to the pool of workers.
func (pool *WorkerPool) AddJob(job Job) {
	go func() { pool.pendingJobs <- job }()
}

// EndJobs tells the Worker Pool that there are no more jobs incoming
// This internally closes the channel for incoming jobs.
func (pool *WorkerPool) EndJobs() {
	close(pool.pendingJobs)
}

// workerRoutine corresponds to the routine in which a worker runs until it is done.
func workerRoutine(pool *WorkerPool) {
	// While there are jobs to process:
	for job := range pool.pendingJobs {
		result := job.Process()
		pool.finishedJobs <- result
	}

	// Mark this worker as finished:
	pool.workersActive.Done()
}

// waitForWorkersRoutine corresponds to a routine that waits until all workers finish.
func waitForWorkersRoutine(pool *WorkerPool) {
	pool.workersActive.Wait()
	close(pool.finishedJobs)
}
