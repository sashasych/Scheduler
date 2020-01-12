package schedule

const (
	DEFAULT_TASK_CAPACITY = 100
)

type Scheduler struct {
	tasks chan Task
}

func NewScheduler(capacity int) *Scheduler {
	return &Scheduler{
		make(chan Task, capacity),
	}
}

func (s *Scheduler) AddTask(task *Task) {
	s.tasks <- *task
}

func (s *Scheduler) Run() {
	go func() {
		for {
			task := <-s.tasks
			go task.start()
		}
	}()
}
