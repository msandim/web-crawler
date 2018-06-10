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
}

// WorkerPool is
type WorkerPool struct {
	nWorkers      int
	pendingJobs   chan Job
	finishedJobs  chan JobResult
	workersActive *sync.WaitGroup
}

var nJobsActive sync.Mutex

// New generates a WorkerPool struct and runs "nWorkers" workers
func New(nWorkers int, jobs chan Job, jobResults chan JobResult) *WorkerPool {

	pool := &WorkerPool{
		nWorkers:      nWorkers,
		pendingJobs:   jobs,
		finishedJobs:  jobResults,
		workersActive: &sync.WaitGroup{},
	}

	return pool
}

// Run initiates the Worker Pool:
func (pool *WorkerPool) Run() {
	// Create a go routine for each worker:
	for i := 0; i < pool.nWorkers; i++ {
		pool.workersActive.Add(1)
		go workerRoutine(pool)
	}

	go waitForWorkersRoutine(pool)
}

// AddJob adds a job to the pool of workers:
func (pool *WorkerPool) AddJob(job Job) {
	select {
	case pool.pendingJobs <- job: // some other worker can do it:

	default: // if the channel is full, do the job synchronously
		result := job.Process()
		pool.finishedJobs <- result
	}
}

// EndJobs tells the Worker Pool that there are no more jobs incoming
// This internally closes the channel for incoming jobs
func (pool *WorkerPool) EndJobs() {
	fmt.Println("fechei pending jobs")
	close(pool.pendingJobs)
}

func workerRoutine(pool *WorkerPool) {
	// While there are jobs to process:
	for job := range pool.pendingJobs {
		// Add more jobs to chan jobs:
		result := job.Process()
		pool.finishedJobs <- result
	}

	pool.workersActive.Done()
	//fmt.Println("finished worker routine")
}

func waitForWorkersRoutine(pool *WorkerPool) {
	//fmt.Println("à espera que os workers terminem")
	pool.workersActive.Wait()
	//fmt.Println("os workers terminaram")
	close(pool.finishedJobs)
}
