package schedule

type PeriodicTaskWithRestrict struct {
}

type Task struct {
	id        string                             // ID задачи
	projectID string                             // ID проекта
	done      chan struct{}                      // Индикатор выполненности задачи
	paused    chan struct{}                      // Индикатор паузы
	stopped   chan struct{}                      // Индикатор остановки
	action    func(chan struct{}, chan struct{}) // Канал паузы, канал остановки
}

func (t *Task) SetAction(action func(chan struct{}, chan struct{})) *Task {
	t.action = action
	return t
}

// Задание id задачи
func (t *Task) SetID(id string) *Task {
	t.id = id
	return t
}

func (t *Task) Done() {
	t.done <- struct{}{}
}

func commonTask() *Task {
	return &Task{
		//days:    make(map[int]struct{}),
		paused:  make(chan struct{}),
		stopped: make(chan struct{}),
		done:    make(chan struct{}),
	}
}

type TasksInterface interface {
	start(stopperFunc func(taskID string))
	getID() string
	Done()
}
