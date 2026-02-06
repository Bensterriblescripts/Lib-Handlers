package time

import "time"

// Formats: https://pkg.go.dev/time#Layout

// GetTime returns the local time formatted as time.Kitchen.
//
// Example:
//
//	now := time.GetTime()
func GetTime() string {
	return time.Now().Local().Format(time.Kitchen)
}
// GetUnixTime returns the current UTC Unix timestamp.
//
// Example:
//
//	epoch := time.GetUnixTime()
func GetUnixTime() int64 {
	return time.Now().UTC().Unix()
}
// GetTimestamp returns the current UTC timestamp string.
//
// Example:
//
//	ts := time.GetTimestamp()
func GetTimestamp() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05.000")
}
// GetFullDateTime returns a short local date/time string.
//
// Example:
//
//	display := time.GetFullDateTime()
func GetFullDateTime() string {
	return time.Now().Local().Format("2/1/06 3:04pm")
}

// GetDay returns the current UTC date as YYYY-M-D.
//
// Example:
//
//	day := time.GetDay()
func GetDay() string {
	return time.Now().UTC().Format("2006-1-2")
}
// GetDateArray returns the local day, month, and year values.
//
// Example:
//
//	parts := time.GetDateArray()
func GetDateArray() []int {
	timenow := time.Now().Local()
	return []int{timenow.Day(), int(timenow.Month()), timenow.Year()}
}
