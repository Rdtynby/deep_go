package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	tasks []Task
}

type SchedulerWithSharing struct {
	schedulers []Scheduler
}

func NewSchedulerWithSharing() SchedulerWithSharing {
	schedulers := []Scheduler{
		NewScheduler(),
		NewScheduler(),
		NewScheduler(),
	}

	return SchedulerWithSharing{
		schedulers: schedulers,
	}
}

func (s *SchedulerWithSharing) AddTask(index int, task Task) {
	s.schedulers[index].AddTask(task)
}

func (s *SchedulerWithSharing) ChangeTaskPriority(index int, taskID int, newPriority int) {
	s.schedulers[index].ChangeTaskPriority(taskID, newPriority)
}

func (s *SchedulerWithSharing) GetTask(index int) (Task, bool) {
	task, ok := s.schedulers[index].GetTask()

	for key := range s.schedulers {
		if ok {
			return task, ok
		}

		task, ok = s.schedulers[key].GetTask()
	}

	return task, ok
}

func NewScheduler() Scheduler {
	return Scheduler{
		tasks: make([]Task, 0),
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.tasks = append(s.tasks, task)
	heapify(s.tasks)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	for key := range s.tasks {
		if s.tasks[key].Identifier == taskID {
			s.tasks[key].Priority = newPriority
		}
	}

	heapify(s.tasks)
}

func (s *Scheduler) GetTask() (Task, bool) {
	if len(s.tasks) == 0 {
		return Task{}, false
	}

	task := s.tasks[0]
	s.tasks = s.tasks[1:]
	heapify(s.tasks)

	return task, true
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task, _ := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task, _ = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)
	task1.Priority = 100

	task, _ = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task, _ = scheduler.GetTask()
	assert.Equal(t, task3, task)
}

func TestTraceSharing(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewSchedulerWithSharing()
	scheduler.AddTask(0, task1)
	scheduler.AddTask(1, task2)
	scheduler.AddTask(2, task3)
	scheduler.AddTask(0, task4)
	scheduler.AddTask(1, task5)

	task, _ := scheduler.GetTask(1)
	assert.Equal(t, task5, task)

	task, _ = scheduler.GetTask(1)
	assert.Equal(t, task2, task)

	scheduler.ChangeTaskPriority(0, 1, 100)
	task1.Priority = 100

	task, _ = scheduler.GetTask(1) // No tasks in scheduler 1, taking from task with the highest priority from scheduler 0
	assert.Equal(t, task1, task)

	task, _ = scheduler.GetTask(1) // No tasks in scheduler 1, taking from task with the highest priority from scheduler 0
	assert.Equal(t, task4, task)

	task, _ = scheduler.GetTask(1) // No tasks in scheduler 1 or 0, taking from task with the highest priority from scheduler 2
	assert.Equal(t, task3, task)
}
