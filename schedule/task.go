package schedule

import (
	"github.com/sashasych/scheduler/util"
	"time"
)

type Task struct {
	timestamp time.Time
	delay     time.Duration
	action    func(t *Task)
	out       chan bool
	done      chan bool
}

func NewTask(timestamp time.Time, action func(t *Task)) (*Task, error) {
	delay, err := util.ComputeDelay(timestamp)
	if err != nil {
		return nil, err
	}
	return &Task{
		timestamp: timestamp,
		delay:     delay,
		action:    action,
		out:       make(chan bool),
		done:      make(chan bool),
	}, nil
}
func (t *Task) Start() {
	for {
		select {
		case <-t.out:
			return
		case <-t.done:
			return
		case <-time.After(t.delay):
			go t.action(t)
		}
	}

}

func (t *Task) Done() {
	t.done <- true
}

func (t *Task) Pause() {

}
