package time

import "time"

// Formats: https://pkg.go.dev/time#Layout

func GetTime() string {
	return time.Now().Local().Format(time.Kitchen)
}
func GetUnixTime() int64 {
	return time.Now().UTC().Unix()
}
func GetTimestamp() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05.000")
}
func GetFullDateTime() string {
	return time.Now().Local().Format("2/1/06 3:04pm")
}

func GetDay() string {
	return time.Now().UTC().Format("2006-1-2")
}
func GetDateArray() []int {
	timenow := time.Now().Local()
	return []int{timenow.Day(), int(timenow.Month()), timenow.Year()}
}
