package schedule

import (
	"fmt"
	"time"

	"github.com/scheduler/util"
)

func Single(targetTime *time.Time) *SingleTask {
	task := commonTask()
	if delay, err := util.ComputeInterval(*targetTime, time.Now()); nil == err {
		return &SingleTask{
			task:  task,
			delay: delay,
		}
	}
	return &SingleTask{
		task: task,
	}
}

type SingleTask struct {
	task  *Task
	delay time.Duration // Задержка перед выполнением задачи
}

func (t *SingleTask) Done() {
	t.task.done <- struct{}{}
}

func (t *SingleTask) SetID(id string) *SingleTask {
	t.task.SetID(id)
	return t
}

func (t *SingleTask) SetAction(action func(chan struct{}, chan struct{})) *SingleTask {
	t.task.SetAction(action)
	return t
}

func (t *SingleTask) getID() string {
	return t.task.id
}

func (t *SingleTask) start(stopperFunc func(taskID string)) {
	var nextInterval time.Duration
	fmt.Println("hello1")
	fmt.Println(t.delay)
	if t.delay <= 0 {
		stopperFunc(t.task.id)
		return
	}
	fmt.Println("hello!")
	nextInterval = t.delay
	for {
		select {
		case <-t.task.done:
			return
		case <-time.After(nextInterval):
			t.task.action(t.task.paused, t.task.stopped)
		}
		t.task.Done()
	}
}
