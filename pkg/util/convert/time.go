package convert

import (
	"strings"
	"time"
)

const timeFormat = time.RFC3339

func TimeToStringFmt(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(timeFormat)
}

func StringToTimeFmt(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(timeFormat, s)
	if err != nil {
		return nil
	}
	return &t
}

// ParseTime returns (zeroTime, false) on parse failure
func ParseTime(raw string, fmt string) (time.Time, bool) {
	if raw = strings.TrimSpace(raw); raw != "" {
		if t, err := time.ParseInLocation(fmt, raw, time.Local); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
