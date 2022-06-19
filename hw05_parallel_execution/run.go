package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(tasks *[]Task, T chan<- error, mu *sync.RWMutex, wg *sync.WaitGroup) {
	var task Task
	pointer := *tasks
	mu.Lock()
	task, *tasks = pointer[0], pointer[1:]
	mu.Unlock()
	wg.Done()
	T <- task()
	if len(*tasks) == 0 {
		close(T)
	}
}

func Run(tasks []Task, n, m int) error {
	var errorsCnt int
	mu := sync.RWMutex{}
	wg := sync.WaitGroup{}

	var T = make(chan error, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(&tasks, T, &mu, &wg)
	}

	for result := range T {
		if result != nil {
			errorsCnt++
		}

		if errorsCnt >= m {
			break
		}

		if len(tasks) > 0 {
			wg.Add(1)
			go worker(&tasks, T, &mu, &wg)
		}
	}

	wg.Wait()
	if errorsCnt >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
