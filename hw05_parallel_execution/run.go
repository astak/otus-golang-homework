package hw05parallelexecution

import (
	"errors"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) (result error) {
	abort := make(chan struct{})

	chTasks := newTasksChannel(tasks, abort)
	chResults := newResultsChannel(chTasks, n)

	totalErrors := 0
	for err := range chResults {
		if err == nil {
			continue
		}

		totalErrors++
		if m > 0 && totalErrors >= m && result == nil {
			close(abort)
			result = ErrErrorsLimitExceeded
		}
	}

	return
}

func newTasksChannel(tasks []Task, abort <-chan struct{}) <-chan Task {
	result := make(chan Task)

	go func() {
		defer close(result)

		for _, task := range tasks {
			select {
			case result <- task:
			case <-abort:
				return
			}
		}
	}()

	return result
}

func newResultsChannel(tasks <-chan Task, workersCount int) <-chan error {
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
				result <- task()
			}
		}()
	}

	return result
}
