package schedule

import (
	"time"

	"github.com/scheduler/util"
)

type TargetTime struct {
	hour     uint8
	minute   uint8
	second   uint8
	timezone int8
}

type PeriodicTaskByDays struct {
	task           *Task
	days           map[int]struct{} // Мапа дней недели в расписании
	targetTime     TargetTime       // Время (час мин сек) для запуска задачи по дням недели
	startTimeDelay time.Duration    // Время задержки перед началом работы правила
	endTimeDelay   time.Duration    // Время оконания работы правила
}

func (t *PeriodicTaskByDays) Done() {
	t.task.done <- struct{}{}
}

func (t *PeriodicTaskByDays) SetID(id string) *PeriodicTaskByDays {
	t.task.SetID(id)
	return t
}

func (t *PeriodicTaskByDays) SetAction(action func(chan struct{}, chan struct{})) *PeriodicTaskByDays {
	t.task.SetAction(action)
	return t
}

func (t *PeriodicTaskByDays) getID() string {
	return t.task.id
}

func commonPeriodicByDaysTask() *PeriodicTaskByDays {
	return &PeriodicTaskByDays{
		days:       make(map[int]struct{}),
		targetTime: TargetTime{},
		task:       commonTask(),
	}
}

func PeriodicByDays(days []int, hour uint8, minute uint8, second uint8, timezone int8, startTime *time.Time, endTime *time.Time, timeNow time.Time) *PeriodicTaskByDays {
	task := commonPeriodicByDaysTask()
	if startTime != nil {
		if startRuleDelay, err := util.ComputeInterval(*startTime, timeNow); nil == err {
			task.startTimeDelay = startRuleDelay
		} else {
			task.startTimeDelay = 0
		}
	} else {
		task.startTimeDelay = 0
	}
	if endTime != nil {
		if endRuleDelay, err := util.ComputeInterval(*endTime, timeNow); nil == err {
			task.endTimeDelay = endRuleDelay
		} else {
			return task
		}
	} else {
		task.endTimeDelay = 0
	}

	task.targetTime.timezone = timezone
	task.targetTime.minute = minute
	task.targetTime.second = second
	hourTimeZone := int8(hour) - timezone
	switch {
	case hourTimeZone < 0 || hourTimeZone > 23:
		task.targetTime.hour = uint8(hourTimeZone % 24)
		for _, day := range days {
			task.days[(day-1)%7] = struct{}{}
		}
	case hourTimeZone > 23:
		task.targetTime.hour = uint8(hourTimeZone % 24)
		for _, day := range days {
			task.days[(day+1)%7] = struct{}{}
		}
	default:
		task.targetTime.hour = uint8(hourTimeZone)
		for _, day := range days {
			task.days[day] = struct{}{}
		}
	}
	return task
}

func (t *PeriodicTaskByDays) getClosestDateByDays(timeNow time.Time, local *time.Location) time.Time {
	dayOfWeek := timeNow.Weekday()
	checkDay := 0
	nextDay := 0
	i := 0
	for ; i < 7; i++ {
		checkDay = (int(dayOfWeek) + i) % 7
		_, ok := t.days[checkDay]
		if ok {

			if checkDay-int(dayOfWeek) == 0 {
				if t.targetTime.hour > uint8(timeNow.Hour()) {
					nextDay = i
					break
				} else if t.targetTime.hour == uint8(timeNow.Hour()) {
					if t.targetTime.minute > uint8(timeNow.Minute()) {
						nextDay = i
						break
					} else if t.targetTime.minute == uint8(timeNow.Minute()) {
						if t.targetTime.second > uint8(timeNow.Second()) {
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
	if i == 7 {
		nextDay = i
	}
	day := nextDay + timeNow.Day()
	closestDate := time.Date(
		timeNow.Year(),
		timeNow.Month(),
		day,
		int(t.targetTime.hour),
		int(t.targetTime.minute),
		int(t.targetTime.second),
		0,
		local,
	)
	return closestDate
}

func (t *PeriodicTaskByDays) start(stopperFunc func(taskID string)) {
	var nextInterval time.Duration
	switch {
	case t.endTimeDelay != 0:
	Loop:
		for {
			select {
			case <-t.task.done:
				return
			case <-time.After(t.startTimeDelay):
				break Loop
			}
		}
		closestDate := t.getClosestDateByDays(time.Now(), time.Local)
		nextInterval, _ = util.ComputeInterval(closestDate, time.Now())
		nextEndDelay := t.endTimeDelay - t.startTimeDelay
		for {
			select {
			case <-t.task.done:
				return
			case <-time.After(nextInterval):
				t.task.action(t.task.paused, t.task.stopped)
			case <-time.After(nextEndDelay):
				stopperFunc(t.task.id)
				return
			}
			nextEndDelay -= nextInterval
			closestDate = t.getClosestDateByDays(time.Now(), time.Local)
			nextInterval, _ = util.ComputeInterval(closestDate, time.Now())
		}
	default:
	Loop1:
		for {
			select {
			case <-t.task.done:
				return
			case <-time.After(t.startTimeDelay):
				break Loop1
			}
		}
		closestDate := t.getClosestDateByDays(time.Now(), time.Local)
		nextInterval, _ = util.ComputeInterval(closestDate, time.Now())
		for {
			select {
			case <-t.task.done:
				return
			case <-time.After(nextInterval):
				t.task.action(t.task.paused, t.task.stopped)
			}
			closestDate = t.getClosestDateByDays(time.Now(), time.Local)
			nextInterval, _ = util.ComputeInterval(closestDate, time.Now())
		}
	}
}
