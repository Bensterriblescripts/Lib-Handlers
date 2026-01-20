package mutate

import (
	"strconv"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

// Converts a slice of int64s to a slice of strings
func ISliceToString(slice []int64) []string {
	stringslice := make([]string, len(slice))
	for index, value := range slice {
		stringslice[index] = strconv.FormatInt(value, 10)
	}
	return stringslice
}

// Converts a slice of strings to a slice of int64s
func SSlicetoISlice(slice []string) []int64 {
	intslice := make([]int64, len(slice))
	for _, value := range slice {
		if intval, err := ErrorExists(strconv.ParseInt(value, 10, 64)); err {
			return nil
		} else {
			intslice = append(intslice, intval)
		}
	}
	return intslice
}
