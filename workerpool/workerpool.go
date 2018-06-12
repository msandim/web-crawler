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
	nWorkers       int
	pendingJobs    chan Job
	pendingAddJobs *sync.WaitGroup
	finishedJobs   chan JobResult
	workersActive  *sync.WaitGroup
}

// New generates a WorkerPool struct and runs "nWorkers" workers.
func New(nWorkers int) *WorkerPool {
	pool := &WorkerPool{
		nWorkers:       nWorkers,
		pendingJobs:    make(chan Job),
		pendingAddJobs: &sync.WaitGroup{},
		finishedJobs:   make(chan JobResult),
		workersActive:  &sync.WaitGroup{},
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

// GetResultsChannel returns the channel from which job results can we collected.
func (pool *WorkerPool) GetResultsChannel() chan JobResult {
	return pool.finishedJobs
}

// AddJob adds a job to the pool of workers.
func (pool *WorkerPool) AddJob(job Job) {
	pool.pendingAddJobs.Add(1)
	go func() {
		pool.pendingJobs <- job
		pool.pendingAddJobs.Done()
	}()
}

// EndJobs tells the Worker Pool that there are no more jobs incoming
// This internally closes the channel for incoming jobs.
// This function call may block if there are jobs waiting to be added to the pendingJobs channel
// (as a result of AddJob()).
func (pool *WorkerPool) EndJobs() {
	pool.pendingAddJobs.Wait()
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
