package utils

import (
	"time"
)

const layout = "2006-01-02 15:04:05"

func ConvertTime(timeString string) time.Time {
	t, err := time.ParseInLocation(layout, timeString, time.Local)
	if err != nil {
		panic(err)
	}

	return t
}
