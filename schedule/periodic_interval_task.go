package schedule

import "time"

type PeriodicTaskByInterval struct {
	task     *Task
	interval time.Duration // Интервал выполнения периодичной задачи
	delay    time.Duration
}

func (t *PeriodicTaskByInterval) Done() {
	t.task.done <- struct{}{}
}

func (t *PeriodicTaskByInterval) SetID(id string) *PeriodicTaskByInterval {
	t.task.SetID(id)
	return t
}

func (t *PeriodicTaskByInterval) SetAction(action func(chan struct{}, chan struct{})) *PeriodicTaskByInterval {
	t.task.SetAction(action)
	return t
}

func (t *PeriodicTaskByInterval) getID() string {
	return t.task.id
}

// Создание периодической задачи с заданным интервалов
func PeriodicByInterval(interval, unit time.Duration) *PeriodicTaskByInterval {
	return &PeriodicTaskByInterval{
		task:     commonTask(),
		interval: interval * unit,
		delay:    0,
	}
}

// Задание задержки перед первым выполнем задачи
func (t *PeriodicTaskByInterval) SetDelay(interval, unit time.Duration) *PeriodicTaskByInterval {
	t.delay = interval * unit
	return t
}

// Функции, создающие периодические задачи по выбранным интервалам
func EveryDay() *PeriodicTaskByInterval {
	return PeriodicByInterval(24, time.Hour)
}

func EveryHour() *PeriodicTaskByInterval {
	return PeriodicByInterval(1, time.Hour)
}

func EveryMinute() *PeriodicTaskByInterval {
	return PeriodicByInterval(1, time.Minute)
}

func EverySecond() *PeriodicTaskByInterval {
	return PeriodicByInterval(1, time.Second)
}

func (t *PeriodicTaskByInterval) start(stopperFunc func(taskID string)) {
Loop1:
	for {
		select {
		case <-t.task.done:
			return
		case <-time.After(t.delay):
			break Loop1
		}
	}
	for {
		select {
		case <-t.task.done:
			return
		case <-time.After(t.interval):
			t.task.action(t.task.paused, t.task.stopped)
		}
	}
}
