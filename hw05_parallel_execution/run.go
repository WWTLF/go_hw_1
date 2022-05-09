package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var errCountMutex sync.RWMutex
	var errCount int32
	taskQue := make(chan Task, len(tasks))
	for _, t := range tasks {
		taskQue <- t
	}
	close(taskQue)

	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(taskQue <-chan Task, mu *sync.RWMutex) {
			defer wg.Done()
			for {
				mu.RLock()
				if errCount >= int32(m) && m > 0 {
					mu.RUnlock()
					return
				}
				mu.RUnlock()

				t, ok := <-taskQue
				if !ok {
					return
				}

				err := t()
				if err != nil {
					mu.Lock()
					errCount++
					mu.Unlock()
				}
			}
		}(taskQue, &errCountMutex)
	}

	wg.Wait()

	if errCount >= int32(m) && m > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
