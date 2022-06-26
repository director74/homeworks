package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	var errorsCnt int
	var errorsChan = make(chan error)
	var done = make(chan struct{})

	go func(done <-chan struct{}) {
		mu := sync.RWMutex{}
		wg := sync.WaitGroup{}
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for len(tasks) > 0 {
					select {
					case <-done:
						return
					default:
						var task Task
						mu.Lock()
						task, tasks = tasks[0], tasks[1:]
						mu.Unlock()

						errorsChan <- task()
					}
				}
			}()
		}
		wg.Wait()
		close(errorsChan)
	}(done)

	for result := range errorsChan {
		if result != nil {
			errorsCnt++
		}

		if errorsCnt == m {
			close(done)
		}
	}

	if errorsCnt >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
