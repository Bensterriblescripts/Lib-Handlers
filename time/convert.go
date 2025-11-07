package time

import "time"

func ConvertUnixToTimestamp(date int64) string {
	return time.Unix(date, 0).UTC().Format("2006-01-02 15:04:05.000")
}
