package schedule

import (
	"sync"
)

type Scheduler struct {
	tasks map[string]*Task
	mu    sync.Mutex
}

func NewScheduler(arrayTasks []*Task) *Scheduler {
	tasksMap := make(map[string]*Task)
	scheduler := Scheduler{
		mu: sync.Mutex{},
	}
	for _, task := range arrayTasks {
		scheduler.mu.Lock()
		tasksMap[task.id] = task
		scheduler.mu.Unlock()
	}
	scheduler.tasks = tasksMap
	return &scheduler
}

func (s *Scheduler) AddTask(task *Task) {
	s.mu.Lock()
	s.tasks[task.id] = task
	s.mu.Unlock()
	go task.start(s.stopperFunc)
}

func (s *Scheduler) stopperFunc(taskID string) {
	s.mu.Lock()
	delete(s.tasks, taskID)
	s.mu.Unlock()
}

func (s *Scheduler) DeleteTask(taskID string) {
	s.mu.Lock()
	s.tasks[taskID].Done()
	delete(s.tasks, taskID)
	s.mu.Unlock()
}

func (s *Scheduler) Run() {
	s.mu.Lock()
	for _, task := range s.tasks {
		go task.start(s.stopperFunc)
	}
	s.mu.Unlock()

}
