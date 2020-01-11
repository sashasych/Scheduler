package schedule

import "time"

const (
	DEFAULT_TASK_CAPACITY = 100
)

//

type Scheduler struct {
	tasks chan Task
}

func NewScheduler(capacity int) *Scheduler {
	return &Scheduler{
		make(chan Task, capacity),
	}
}

//func (s *Scheduler) StartTaskReceiver() {
//	go func() {
//		for {
//			s.ScheduleTask()
//		}
//	}()
//}

//func (s *Scheduler) ScheduleTask() {
//	go func() {
//		for {
//			// TODO - распараллелить чтение из канала
//			task := <- s.tasks
//			timeout := time.Duration(task.sleepTime)
//			select {
//			case <- time.After(timeout * time.Second):
//				fmt.Println(task.message)
//			}
//		}
//	}()
//}

func (s *Scheduler) AddTask(timestamp time.Time, action func(task *Task)) {
	if task, err := NewTask(timestamp, action); err == nil {
		s.tasks <- *task
	}
}
func (s *Scheduler) Run() {
	go func() {
		for {
			task := <-s.tasks
			go task.Start()
		}
	}()
}
