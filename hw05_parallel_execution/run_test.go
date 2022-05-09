package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				MySleep(t, taskSleep)
				// time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestErrCountLogic(t *testing.T) {
	t.Run("Error count m > 0 error", func(t *testing.T) {
		tasks := make([]Task, 0, 5)
		tasks = append(tasks, func() error {
			return nil
		})
		tasks = append(tasks, func() error {
			return nil
		})
		tasks = append(tasks, func() error {
			return errors.New("I am not OK")
		})
		tasks = append(tasks, func() error {
			return nil
		})
		tasks = append(tasks, func() error {
			return nil
		})
		tasks = append(tasks, func() error {
			return errors.New("I am not OK")
		})

		err := Run(tasks, 2, 2)
		require.Error(t, err, ErrErrorsLimitExceeded)
	})

	t.Run("Error count m == 0 no error", func(t *testing.T) {
		tasks := make([]Task, 0, 5)
		tasks = append(tasks, func() error {
			return nil
		})
		tasks = append(tasks, func() error {
			return nil
		})
		tasks = append(tasks, func() error {
			return errors.New("I am not OK")
		})
		tasks = append(tasks, func() error {
			return nil
		})
		tasks = append(tasks, func() error {
			return nil
		})
		tasks = append(tasks, func() error {
			return errors.New("I am not OK")
		})

		err := Run(tasks, 2, 0)
		require.NoError(t, err)
	})

	t.Run("Error count m = 3 no error", func(t *testing.T) {
		tasks := make([]Task, 0, 5)
		tasks = append(tasks, func() error {
			return nil
		})
		tasks = append(tasks, func() error {
			return nil
		})
		tasks = append(tasks, func() error {
			return errors.New("I am not OK")
		})
		tasks = append(tasks, func() error {
			return nil
		})
		tasks = append(tasks, func() error {
			return nil
		})
		tasks = append(tasks, func() error {
			return errors.New("I am not OK")
		})

		err := Run(tasks, 2, 3)
		require.NoError(t, err)
	})
}

func MySleep(t *testing.T, d time.Duration) {
	t.Helper()
	if d == 0 {
		return
	}
	start := time.Now()
	require.Eventually(t, func() bool {
		ellapsedTime := time.Since(start)
		return ellapsedTime >= d
	}, d*10, time.Millisecond)
}
