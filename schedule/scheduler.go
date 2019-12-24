package schedule

import (
	"fmt"
	"time"
)

const (
	DEFAULT_TASK_CAPACITY = 100
)

type Task struct {
	sleepTime int
	message string
}

type Scheduler struct {
	tasks chan Task
}

func NewScheduler(capacity int) *Scheduler {
	return &Scheduler{
		make(chan Task, capacity),
	}
}

func (s *Scheduler) StartTaskReceiver() {
	go func() {
		for {
			s.ScheduleTask()
		}
	}()
}

func (s *Scheduler) ScheduleTask() {
	go func() {
		for {
			// TODO - распараллелить чтение из канала
			task := <- s.tasks
			timeout := time.Duration(task.sleepTime)
			select {
			case <- time.After(timeout * time.Second):
				fmt.Println(task.message)
			}
		}
	}()
}

func (s *Scheduler) AddTask(timeout int, message string) {
	s.tasks <- Task{timeout, message}
}