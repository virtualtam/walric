package formatter

import "time"

// FormatDateAsUTC returns a string representation of a given date, in the UTC
// standard.
func FormatDateAsUTC(date time.Time) string {
	return date.UTC().Format("2006-01-02 15:04:05 MST")
}
