package time

// UnixSince reports whether the target time is more than the provided seconds ago.
//
// Example:
//
//	expired := time.UnixSince(lastSeen, 300)
func UnixSince(targettime int64, secondssince int64) bool {
	return GetUnixTime()-targettime > secondssince
}
// UnixUntil reports whether the target time is within the provided seconds from now.
//
// Example:
//
//	soon := time.UnixUntil(eventTime, 60)
func UnixUntil(targettime int64, secondssince int64) bool {
	return GetUnixTime()-targettime < secondssince
}
