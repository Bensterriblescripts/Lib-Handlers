package time

func UnixSince(targettime int64, secondssince int64) bool {
	return GetUnixTime()-targettime > secondssince
}
func UnixUntil(targettime int64, secondssince int64) bool {
	return GetUnixTime()-targettime < secondssince
}
