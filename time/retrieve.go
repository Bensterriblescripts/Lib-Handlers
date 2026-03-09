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

// GetPastTimestamp returns the UTC timestamp string minus the number of minutes specified.
//
// Example:
//
//	ts := time.GetPastTimestamp(10) // 10 minutes ago
func GetPastTimestamp(minutes int) string {
	return time.Now().UTC().Add(-time.Minute * time.Duration(minutes)).Format("2006-01-02 15:04:05.000")
}

//
// GetPastTimestampSeconds returns the UTC timestamp string minus the number of seconds specified.
//
// Example:
//
//	ts := time.GetPastTimestampSeconds(10) // 10 seconds ago
func GetPastTimestampSeconds(seconds int) string {
	return time.Now().UTC().Add(-time.Second * time.Duration(seconds)).Format("2006-01-02 15:04:05.000")
}

//
// GetFutureTimestamp returns the UTC timestamp string plus the number of minutes specified.
//
// Example:
//
//	ts := time.GetFutureTimestamp(10) // 10 minutes from now
func GetFutureTimestamp(minutes int) string {
	return time.Now().UTC().Add(time.Minute * time.Duration(minutes)).Format("2006-01-02 15:04:05.000")
}

//
// GetFutureTimestampSeconds returns the UTC timestamp string plus the number of seconds specified.
//
// Example:
//
//	ts := time.GetFutureTimestampSeconds(10) // 10 seconds from now
func GetFutureTimestampSeconds(seconds int) string {
	return time.Now().UTC().Add(time.Second * time.Duration(seconds)).Format("2006-01-02 15:04:05.000")
}

// GetFullDateTime returns a short local date/time string.
//
// Example:
//
//	display := time.GetFullDateTime()
func GetFullDayTime() string {
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
