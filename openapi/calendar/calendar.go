package calendar

import (
	"fmt"
	"time"
)

const iso8601 = "2006-01-02T15:04:05+09:00"

type Calendar struct {
	endPoint string
}

func NewDefaultCalendar() *Calendar {
	return &Calendar{
		endPoint: "https://api.dooray.com",
	}
}
func NewCalendar(endPoint string) *Calendar {
	return &Calendar{
		endPoint: endPoint,
	}
}

func ConvertISO8601ToTime(t string) (time.Time, error) {
	parsedTime, err := time.Parse(iso8601, t)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}
	return parsedTime, nil
}

func FormatTimeToISO8601(t time.Time) string {

	formattedTime := t.Format(iso8601)
	return formattedTime
}
