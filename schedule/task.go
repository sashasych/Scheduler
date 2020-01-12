package schedule

import (
	"fmt"
	"github.com/sashasych/scheduler/util"
	"time"
)

type Task struct {
	interval time.Duration              // Интервал выполнения периодичной задачи
	delay    time.Duration              // Задержка перед первым выполнением задачи
	action   func(chan bool, chan bool) // Канал паузы, канал остановки
	repeat   bool                       // Является ли задачи периодической
	done     chan bool                  // Индикатор выполненности задачи
	paused   chan bool                  // Индикатор паузы
	stopped  chan bool                  // Индикатор остановки
}

// Создание задачи, которую необходимо выполнить единожды
func Single(targetTime time.Time) *Task {
	task := commonTask()
	task.repeat = false
	if delay, err := util.ComputeInterval(targetTime); nil == err {
		task.delay = delay
	} else {
		fmt.Println("Delay in past, task will execute immediately")
	}
	return task
}

// Создание периодической задачи с заданным интервалов
func Periodic(interval, unit time.Duration) *Task {
	task := commonTask()
	task.interval = interval * unit
	task.repeat = true
	return task
}

// Задание задержки перед первым выполнем задачи
func (t *Task) SetDelay(interval, unit time.Duration) *Task {
	if t.repeat {
		t.delay = interval * unit
	}
	return t
}

func (t *Task) SetAction(action func(chan bool, chan bool)) *Task {
	t.action = action
	return t
}

// Запуск задачи
// Если задача должна выполниться единожды, то выполняем ее после подсчитанной задержки
// Если задача периодическая, выполняем ее через одинаковые интервалы времени + если она имеет задержку
// логика выполнения первой итерации совпадает с логикой выполнения задачи "на один раз"
func (t *Task) start() {
	nextInterval := t.delay
	for {
		select {
		case <-t.done:
			return
		case <-time.After(nextInterval):
			t.action(t.paused, t.stopped)
		}
		nextInterval = t.interval
		if !t.repeat {
			t.Done()
		}
	}
}

func (t *Task) Done() {
	t.done <- true
}

func (t *Task) Pause() {
	t.Done()
	t.delay = 0
}

func (t *Task) Resume() {
	t.done = make(chan bool)
	t.start()
}

// Функции, создающие периодические задачи по выбранным интервалам
func EveryDay() *Task {
	return Periodic(24, time.Hour)
}

func EveryHour() *Task {
	return Periodic(1, time.Hour)
}

func EveryMinute() *Task {
	return Periodic(1, time.Minute)
}

func EverySecond() *Task {
	return Periodic(1, time.Second)
}

func commonTask() *Task {
	return &Task{
		paused:  make(chan bool),
		stopped: make(chan bool),
		done:    make(chan bool),
	}
}
