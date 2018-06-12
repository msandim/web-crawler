package workerpool

import (
	"testing"
	"time"
)

type testJob struct {
	id int
}

type testJobResult struct {
	job *testJob
}

func (job *testJob) Process() JobResult {
	return &testJobResult{job: job}
}

func (result *testJobResult) GetJob() Job {
	return result.job
}

func TestNew(t *testing.T) {
	workerPool := New(10)

	if workerPool.nWorkers != 10 {
		t.Errorf("Number of workers was incorrect, got: %d, want: %d.", workerPool.nWorkers, 10)
	}

	if workerPool.pendingJobs == nil {
		t.Errorf("Incorrect pendingJobs channel initialized.")
	}

	if workerPool.finishedJobs == nil {
		t.Errorf("Incorrect finishedJobs channel initialized.")
	}

	if workerPool.pendingAddJobs == nil {
		t.Errorf("Sync variable pendingAddJobs was not initialized.")
	}

	if workerPool.workersActive == nil {
		t.Errorf("Sync variable workersActive was not initialized.")
	}
}

func TestRun(t *testing.T) {
	workerPool := New(10)
	results := workerPool.GetResultsChannel()

	go workerPool.Run()
	workerPool.AddJob(&testJob{id: 1})
	workerPool.AddJob(&testJob{id: 2})
	workerPool.AddJob(&testJob{id: 3})
	workerPool.EndJobs()
	jobsDone := 0
	job1 := 0
	job2 := 0
	job3 := 0

	for result := range results {
		job := result.GetJob().(*testJob)
		jobsDone++

		switch job.id {
		case 1:
			job1++
		case 2:
			job2++
		case 3:
			job3++
		default:
		}
	}

	if jobsDone != 3 {
		t.Errorf("Number of jobs done was incorrect, got: %d, want: %d.", jobsDone, 3)
	}
	if job1 != 1 {
		t.Errorf("Number of job1 executions was wrong, got: %d, want: %d.", job1, 1)
	}
	if job2 != 1 {
		t.Errorf("Number of job2 executions was wrong, got: %d, want: %d.", job2, 1)
	}
	if job3 != 1 {
		t.Errorf("Number of job3 executions was wrong, got: %d, want: %d.", job3, 1)
	}
}

func TestRun2(t *testing.T) {
	workerPool := New(10)
	results := workerPool.GetResultsChannel()
	go workerPool.Run()
	workerPool.AddJob(&testJob{id: 1})

	_, ok := <-results
	if !ok {
		t.Errorf("Results channel was closed without EndJobs() call.")
	}

	workerPool.EndJobs()
	time.Sleep(1 * time.Second) // wait for routines of the workerpool to close the channels

	_, ok = <-results
	if ok {
		t.Errorf("Results channel still open after EndJobs() call.")
	}
}
