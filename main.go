package main

import (
	"fmt"
	"github.com/sashasych/scheduler/schedule"
	"time"
)

func main() {
	scheduler := schedule.NewScheduler(schedule.DEFAULT_TASK_CAPACITY)
	scheduler.Run()
	scheduler.AddTask(  time.Date(2020, 1, 11, 21, 52, 0, 0, time.UTC), func(task *schedule.Task) {
		fmt.Println("hello!")
		task.Done()
	})
	scheduler.AddTask(  time.Date(2020, 1, 11, 21, 52, 5, 0, time.UTC), func(task *schedule.Task) {
		fmt.Println("hello!!!")
		task.Done()
	})
	time.Sleep(time.Second * 100)
	//scheduler.AddTask(1578775800, "after 5 secs")
}