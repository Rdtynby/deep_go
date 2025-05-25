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

type Task struct {
	job     func()
	counter int
}

type WorkerPool struct {
	tasks  chan Task
	wg     sync.WaitGroup
	closed bool
}

func NewWorkerPool(workersNumber int) *WorkerPool {
	wp := &WorkerPool{
		tasks:  make(chan Task, workersNumber*2),
		closed: false,
	}

	wp.wg.Add(workersNumber)

	for i := 0; i < workersNumber; i++ {
		go wp.PerformTasks()
	}

	return wp
}

// Return an error if the pool is full
func (wp *WorkerPool) AddTask(task func(), counter int) error {
	if wp.closed {
		return fmt.Errorf("already closed")
	}

	multitask := Task{
		job:     task,
		counter: counter,
	}

	select {
	case wp.tasks <- multitask:
		return nil
	default:
		return fmt.Errorf("task buffer is full")
	}
}

// Shutdown all workers and wait for all
// tasks in the pool to complete
func (wp *WorkerPool) Shutdown() {
	if wp.closed {
		return
	}

	close(wp.tasks)
	wp.wg.Wait()
	wp.closed = true
}

func (wp *WorkerPool) PerformTasks() {
	defer wp.wg.Done()

	for task := range wp.tasks {
		for i := 0; i < task.counter; i++ {
			task.job()
		}
	}
}

func TestWorkerPool(t *testing.T) {
	var counter atomic.Int32
	task := func() {
		time.Sleep(time.Millisecond * 500)
		counter.Add(1)
	}

	pool := NewWorkerPool(2)
	_ = pool.AddTask(task, 1)
	_ = pool.AddTask(task, 1)
	_ = pool.AddTask(task, 1)

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(2), counter.Load())

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(3), counter.Load())

	_ = pool.AddTask(task, 1)
	_ = pool.AddTask(task, 1)
	_ = pool.AddTask(task, 3)
	pool.Shutdown() // wait tasks

	assert.Equal(t, int32(8), counter.Load())

	err := pool.AddTask(task, 1)
	assert.EqualError(t, err, "already closed")
}
