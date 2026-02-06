package time

import "time"

// ConvertUnixToTimestamp formats a Unix timestamp into a UTC date string.
//
// Example:
//
//	text := time.ConvertUnixToTimestamp(1700000000)
func ConvertUnixToTimestamp(date int64) string {
	return time.Unix(date, 0).UTC().Format("2006-01-02 15:04:05.000")
}
