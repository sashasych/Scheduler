package util

import (
	"errors"
	"time"
)

// Рассчет времени ожидания исходя из текущей метки времени
func ComputeInterval(timestamp time.Time, timeNow time.Time) (time.Duration, error) {
	result := timestamp.Sub(timeNow)
	if result < 0 {
		return 0, errors.New("time had past")
	}
	return result, nil
}
