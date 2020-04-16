package schedule

import (
	"fmt"
	"sync"
)

type Scheduler struct {
	tasks map[string]TasksInterface
	mu    sync.Mutex
}

// Создание сервиса расписаний
func NewScheduler(arrayTasks []TasksInterface) *Scheduler {
	tasksMap := make(map[string]TasksInterface)
	for _, task := range arrayTasks {
		tasksMap[task.getID()] = task
	}
	return &Scheduler{
		tasks: tasksMap,
		mu:    sync.Mutex{},
	}
}

func (s *Scheduler) AddTask(task TasksInterface) {
	s.mu.Lock()
	id := task.getID()
	if _, ok := s.tasks[id]; ok {
		s.stopAndDeleteTask(id)
	}
	s.tasks[id] = task
	fmt.Println("start task")
	go task.start(s.stopperFunc)
	s.mu.Unlock()
}

// В функцию start у раписания(задачи) передается функция stopper,
func (s *Scheduler) stopperFunc(taskID string) {
	s.mu.Lock()
	delete(s.tasks, taskID)
	s.mu.Unlock()
}

func (s *Scheduler) stopAndDeleteTask(taskID string) {
	s.tasks[taskID].Done()
	delete(s.tasks, taskID)
}

// Единичное удаление расписания(задачи)
func (s *Scheduler) DeleteTask(taskID string) {
	s.mu.Lock()
	if _, ok := s.tasks[taskID]; ok {
		s.tasks[taskID].Done()
		delete(s.tasks, taskID)
	}
	s.mu.Unlock()
}

// Каскадное удаление расписаний(задач)
func (s *Scheduler) DeleteTasks(tasksID []string) {
	for _, taskID := range tasksID {
		s.DeleteTask(taskID)
	}
}

// Запуск расписаний(задач), которые были проинициализированы
func (s *Scheduler) Run() {
	s.mu.Lock()
	for _, task := range s.tasks {
		go task.start(s.stopperFunc)
	}
	s.mu.Unlock()
}
