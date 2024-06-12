package main

import (
	"fmt"
	"time"

	"github.com/scheduler/schedule"
)

func main() {

	// some test comment

	// Инициализируем задачи(пустым массивом задач)
	scheduler := schedule.NewScheduler(nil)
	// Запускаем задачи
	scheduler.Run()

	//Создаем пару задач
	var timeExecute time.Time
	timeExecute = time.Date(2020, 4, 10, 13, 16, 0, 0, time.Local)

	singleTask := schedule.Single(&timeExecute).
		SetAction(func(chan struct{}, chan struct{}) {
			fmt.Println("task #1")
		}).SetID("1")
	periodicTaskByDays := schedule.PeriodicByDays([]int{1, 2, 3, 4, 5, 6}, 14, 40, 0, 0, nil, nil, time.Now()).
		SetAction(func(chan struct{}, chan struct{}) {
			fmt.Println("periodic task")
		}).SetID("2")
	periodicTaskByInterval := schedule.PeriodicByInterval(5, time.Second).SetDelay(10, time.Second).
		SetAction(func(chan struct{}, chan struct{}) {
			fmt.Println("task #1")
		}).SetID("3")
	//Добавляем задачи в сервис расписаний(scheduler)
	scheduler.AddTask(singleTask)
	scheduler.AddTask(periodicTaskByDays)
	scheduler.AddTask(periodicTaskByInterval)

	//time.AfterFunc(10*time.Second, func() {
	//	periodicTask.Pause()
	//})
	//time.AfterFunc(15*time.Second, func() {
	//	periodicTask.Resume()
	//})
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
