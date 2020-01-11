package util

import (
	"errors"
	"time"
)

func ComputeDelay(timestamp time.Time) (time.Duration, error) {
	result := timestamp.Sub(time.Now())
	if result < 0 {
		return 0, errors.New("time had past")
	}
	return result, nil
}