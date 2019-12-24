package main

import (
	"scheduler/schedule"
)

func main() {
	stopped := make(chan bool)
	scheduler := schedule.NewScheduler(schedule.DEFAULT_TASK_CAPACITY)
	scheduler.StartTaskReceiver()
	scheduler.AddTask(2, "after 2 secs")
	scheduler.AddTask(5, "after 5 secs")
	<-stopped
}