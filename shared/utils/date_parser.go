package utils

import "time"

// ParseFlexibleDate accepts "YYYY-MM-DD" or RFC3339 ("2006-01-02T15:04:05Z07:00").
func ParseFlexibleDate(s string) (time.Time, error) {
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339, s)
}
