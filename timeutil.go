package vokiri

import (
	"time"
)

const (
	WAIT_VERYSHORT = 20 * time.Millisecond  //ms
	WAIT_SHORT     = 200 * time.Millisecond //ms
	WAIT_LONG      = 1 * time.Second
)

func Wait(timeout, wait time.Duration, f func() bool) bool {
	for {
		done := f()
		if done {
			return true
		}
		if timeout < 0 {
			break
		}

		time.Sleep(wait)

		timeout -= wait
	}

	return false
}
