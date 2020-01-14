package main

import (
	"fmt"
	"github.com/sashasych/scheduler/schedule"
	"time"
)

func main() {
	scheduler := schedule.NewScheduler(schedule.DEFAULT_TASK_CAPACITY)
	scheduler.Run()

	task := schedule.Single(time.Date(2020, 1, 12, 19, 37, 0, 0, time.UTC)).
		SetAction(func(chan bool, chan bool) {
			fmt.Println("task #1")
		})
	periodicTask := schedule.EverySecond().
		SetAction(func(chan bool, chan bool) {
			fmt.Println("periodic task")
		}).
		SetDelay(1, time.Second)
	scheduler.AddTask(task)
	scheduler.AddTask(periodicTask)
	time.AfterFunc(10 * time.Second, func() {
		periodicTask.Pause()
	})
	time.AfterFunc(15 * time.Second, func() {
		periodicTask.Resume()
	})
	//scheduler.AddTask(  time.Date(2020, 1, 12, 18, 18, 0, 0, time.UTC), func(task *schedule.Task) {
	//	fmt.Println("hello!")
	//	task.Done()
	//})
	//scheduler.AddTask(  time.Date(2020, 1, 12, 18, 18, 10, 0, time.UTC), func(task *schedule.Task) {
	//	fmt.Println("hello!!!")
	//	task.Done()
	//})
	time.Sleep(time.Second * 100)
	//scheduler.AddTask(1578775800, "after 5 secs")
}
