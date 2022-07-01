package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	var errorsCnt int
	errorsChan := make(chan error)
	done := make(chan struct{})
	mu := sync.RWMutex{}
	wg := sync.WaitGroup{}

	worker := func() {
		var task Task
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				mu.Lock()
				leftCnt := len(tasks)
				if leftCnt > 0 {
					task, tasks = tasks[0], tasks[1:]
				}
				mu.Unlock()
				if leftCnt == 0 {
					return
				}
				errorsChan <- task()
			}
		}
	}

	go func() {
		for i := 0; i < n; i++ {
			wg.Add(1)
			go worker()
		}
		wg.Wait()
		close(errorsChan)
	}()

	for result := range errorsChan {
		if result != nil {
			errorsCnt++
		}

		if errorsCnt == m && m > 0 {
			close(done)
		}
	}

	if errorsCnt >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
