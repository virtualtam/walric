package formatter

import "time"

func FormatDateAsUTC(date time.Time) string {
	return date.UTC().Format("2006-01-02 15:04:05 MST")
}
