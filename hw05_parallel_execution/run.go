package hw05parallelexecution

import (
	"errors"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type worker struct {
	tasks       <-chan Task
	done        <-chan struct{}
	results     chan<- error
	jobsRunning int32
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	numTasks := len(tasks)
	chTasks := make(chan Task, numTasks)
	chResults := make(chan error, numTasks)
	chDone := make(chan struct{})
	worker := newWorker(chTasks, chDone, chResults)
	worker.Run(n)

	defer close(chDone)

	for _, t := range tasks {
		chTasks <- t
	}
	close(chTasks)

	totalErrors := 0
	for err := range chResults {
		if err != nil {
			totalErrors++
		}

		if m > 0 && totalErrors >= m {
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}

func newWorker(jobs <-chan Task, done <-chan struct{}, results chan<- error) *worker {
	return &worker{
		tasks:   jobs,
		done:    done,
		results: results,
	}
}

func (w *worker) Run(numJobs int) {
	for i := 0; i < numJobs; i++ {
		go w.doWork()
		atomic.AddInt32(&w.jobsRunning, 1)
	}
}

func (w *worker) doWork() {
	defer func() {
		if remaining := atomic.AddInt32(&w.jobsRunning, -1); remaining == 0 {
			close(w.results)
		}
	}()

	for t := range w.tasks {
		select {
		case <-w.done:
			return
		default:
			w.results <- t()
		}
	}
}
