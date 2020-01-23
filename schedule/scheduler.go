package schedule
type Scheduler struct {
	tasks map[string]*Task
}

func NewScheduler(arrayTasks []*Task) *Scheduler {
	tasksMap := make(map[string]*Task)
	for _, task := range arrayTasks {
		tasksMap[task.id] = task
	}
	return &Scheduler{
		tasks: tasksMap,
	}
}

func (s *Scheduler) AddTask(task *Task) {
	//TODO: добавление задачи в мапу
	go task.start()
}

func (s *Scheduler) Run() {
	for _, task := range s.tasks {
		go task.start()
	}
}
