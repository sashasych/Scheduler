package schedule

import (
	"github.com/sashasych/scheduler/util"
	"time"
)

type TargetTime struct {
	hour   uint8
	minute uint8
	second uint8
}

type Task struct {
	id         string       // ID задачи
	days       map[int]bool // Мапа дней недели в расписании
	targetTime TargetTime
	interval   time.Duration              // Интервал выполнения периодичной задачи
	delay      time.Duration              // Задержка перед первым выполнением задачи
	action     func(chan bool, chan bool) // Канал паузы, канал остановки
	repeat     bool                       // Является ли задачи периодической
	done       chan bool                  // Индикатор выполненности задачи
	paused     chan bool                  // Индикатор паузы
	stopped    chan bool                  // Индикатор остановки
}

// Создание задачи, которую необходимо выполнить единожды
func Single(targetTime time.Time) *Task {
	task := commonTask()
	task.repeat = false
	if delay, err := util.ComputeInterval(targetTime); nil == err {
		task.delay = delay
	} else {
		return task
		//fmt.Println("Delay in past, task will execute immediately")
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

func PeriodicByDays(days []int, hour uint8, minute uint8, second uint8) *Task {
	task := commonTask()
	for _, day := range days {
		task.days[day] = true
	}
	task.targetTime.hour = hour
	task.targetTime.minute = minute
	task.targetTime.second = second
	task.repeat = true
	return task
}

func (t *Task) calculateNextIntervalByDays() time.Duration {
	dayOfWeek := time.Now().Weekday()
	checkDay := 0
	nextDay := 0
	for i := 0; i < 7; i++ {
		checkDay = (int(dayOfWeek) + i) % 7
		_, ok := t.days[checkDay]
		if ok {
			if checkDay-int(dayOfWeek) == 0 {
				if t.targetTime.hour > uint8(time.Now().Hour()) {
					nextDay = i
					break
				} else if t.targetTime.hour == uint8(time.Now().Hour()) {
					if t.targetTime.minute > uint8(time.Now().Minute()) {
						nextDay = i
						break
					} else if t.targetTime.minute == uint8(time.Now().Minute()) {
						if t.targetTime.second > uint8(time.Now().Second()) {
							nextDay = i
							break
						}
					}
				}
			} else {
				nextDay = i
				break
			}
		}
	}
	day := nextDay + time.Now().Day()
	interval, _ := util.ComputeInterval(
		time.Date(
			time.Now().Year(),
			time.Now().Month(),
			day,
			int(t.targetTime.hour),
			int(t.targetTime.minute),
			int(t.targetTime.second),
			0,
			time.Local,
		))
	return interval
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

func (t *Task) SetID(id string) *Task {
	t.id = id
	return t
}

// Запуск задачи
// Если задача должна выполниться единожды, то выполняем ее после подсчитанной задержки
// Если задача периодическая, выполняем ее через одинаковые интервалы времени + если она имеет задержку
// логика выполнения первой итерации совпадает с логикой выполнения задачи "на один раз"
func (t *Task) start(stopperFunc func(taskID string)) {
	var nextInterval time.Duration
	if !t.repeat {
		if t.delay <= 0 {
			stopperFunc(t.id)
			return
		}
		nextInterval = t.delay
	} else {
		nextInterval = t.calculateNextIntervalByDays()
	}
	for {
		select {
		case <-t.done:
			stopperFunc(t.id)
			return
		case <-time.After(nextInterval):
			t.action(t.paused, t.stopped)
		}
		if !t.repeat {
			t.Done()
			stopperFunc(t.id)
		} else {
			nextInterval = t.calculateNextIntervalByDays()
		}
	}
}

func (t *Task) Done() {
	t.done <- true
}

/*func (t *Task) Pause() {
	t.Done()
	t.delay = 0
}*/

/*func (t *Task) Resume() {
	close(t.done)
	t.done = make(chan bool)
	t.start()
}*/

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
		days:    make(map[int]bool),
		paused:  make(chan bool),
		stopped: make(chan bool),
		done:    make(chan bool),
	}
}