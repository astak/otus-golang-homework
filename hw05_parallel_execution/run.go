package hw05parallelexecution

import (
	"errors"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	done := make(chan struct{})
	defer close(done)

	chTasks := newTasksChannel(tasks, done)
	chResults := newResultsChannel(chTasks, n, done)

	totalErrors := 0
	for err := range chResults {
		if err == nil {
			continue
		}

		totalErrors++
		if m > 0 && totalErrors >= m {
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}

func newTasksChannel(tasks []Task, done <-chan struct{}) <-chan Task {
	result := make(chan Task)

	go func() {
		defer close(result)

		for _, task := range tasks {
			select {
			case result <- task:
			case <-done:
				return
			}
		}
	}()

	return result
}

func newResultsChannel(tasks <-chan Task, workersCount int, done <-chan struct{}) <-chan error {
	result := make(chan error)
	workersRunning := int32(workersCount)

	for i := 0; i < workersCount; i++ {
		go func() {
			defer func() {
				if currentWorkers := atomic.AddInt32(&workersRunning, -1); currentWorkers == 0 {
					close(result)
				}
			}()

			for task := range tasks {
				select {
				case result <- task():
				case <-done:
					return
				}
			}
		}()
	}

	return result
}
