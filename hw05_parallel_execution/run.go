package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(tasks *[]Task, t chan<- error, mu *sync.RWMutex, wg *sync.WaitGroup) {
	var task Task
	pointer := *tasks
	mu.Lock()
	task, *tasks = pointer[0], pointer[1:]
	mu.Unlock()
	wg.Done()
	t <- task()
	if len(*tasks) == 0 {
		close(t)
	}
}

func Run(tasks []Task, n, m int) error {
	var errorsCnt int
	mu := sync.RWMutex{}
	wg := sync.WaitGroup{}

	t := make(chan error, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(&tasks, t, &mu, &wg)
	}

	for result := range t {
		if result != nil {
			errorsCnt++
		}

		if errorsCnt >= m {
			break
		}

		if len(tasks) > 0 {
			wg.Add(1)
			go worker(&tasks, t, &mu, &wg)
		}
	}

	wg.Wait()
	if errorsCnt >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
