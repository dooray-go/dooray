package utils

import (
	"fmt"
	"strings"
	"time"
)

type JsonTime time.Time

func (jt *JsonTime) String() string {
	t := time.Time(*jt)
	return FormatTimeToISO8601(t)
}

func (jt JsonTime) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"%s"`, jt.String()), nil
}

func (jt *JsonTime) UnmarshalJSON(b []byte) error {
	timeString := strings.Trim(string(b), `"`)
	time, err := ConvertISO8601ToTime(timeString)
	if err != nil {
		return fmt.Errorf("invalid date format: %s, %v", timeString, err)
	}
	*jt = JsonTime(time)
	return nil
}

func (jt *JsonTime) ToTime() time.Time {
	return time.Time(*jt)
}
