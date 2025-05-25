package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type WorkerPool struct {
	jobs   chan func()
	wg     sync.WaitGroup
	closed bool
}

func NewWorkerPool(workersNumber int) *WorkerPool {
	wp := &WorkerPool{
		jobs:   make(chan func(), workersNumber*2),
		closed: false,
	}

	wp.wg.Add(workersNumber)

	for i := 0; i < workersNumber; i++ {
		go wp.PerformJobs()
	}

	return wp
}

// Return an error if the pool is full
func (wp *WorkerPool) AddTask(task func()) error {
	if wp.closed {
		return fmt.Errorf("already closed")
	}

	select {
	case wp.jobs <- task:
		return nil
	default:
		return fmt.Errorf("job buffer is full")
	}
}

// Shutdown all workers and wait for all
// tasks in the pool to complete
func (wp *WorkerPool) Shutdown() {
	if wp.closed {
		return
	}

	close(wp.jobs)
	wp.wg.Wait()
	wp.closed = true
}

func (wp *WorkerPool) PerformJobs() {
	defer wp.wg.Done()

	for job := range wp.jobs {
		job()
	}
}

func TestWorkerPool(t *testing.T) {
	var counter atomic.Int32
	task := func() {
		time.Sleep(time.Millisecond * 500)
		counter.Add(1)
	}

	pool := NewWorkerPool(2)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(2), counter.Load())

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(3), counter.Load())

	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	pool.Shutdown() // wait tasks

	assert.Equal(t, int32(6), counter.Load())

	err := pool.AddTask(task)
	assert.EqualError(t, err, "already closed")
}
