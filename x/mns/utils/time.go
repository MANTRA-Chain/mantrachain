package utils

import (
	"time"
)

// SecondsToTime converts unix seconds to time
func SecondsToTime(seconds int64) time.Time {
	return time.Unix(seconds, 0)
}

// TimeToSeconds converts a time.Time to unix seconds timestamp
func TimeToSeconds(t time.Time) int64 {
	return t.Unix()
}
