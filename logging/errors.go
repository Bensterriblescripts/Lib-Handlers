package logging

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime/debug"
)

// Error Only
func PanicErr(err error) {
	if err != nil {
		Panic(err.Error())
	}
}
func PrintErr(err error) {
	if err != nil {
		ErrorLog(err.Error())
	}
}
func ErrExists(err error) bool {
	if err != nil {
		ErrorLog(err.Error())
		return true
	}
	return false
}
func RetrieveErr(err error) error {
	if err != nil {
		ErrorLog(err.Error())
		return err
	}
	return nil
}

// Value + Error
func PanicError[T any](value T, err error) T {
	if err != nil {
		Panic(err.Error())
	}
	return value
}
func PrintError[T any](value T, err error) T {
	if err != nil {
		ErrorLog(err.Error())
	}
	return value
}
func ErrorExists[T any](value T, err error) (T, bool) {
	if err != nil {
		ErrorLog(err.Error())
		return value, true
	}
	return value, false
}
func RetrieveValue[T any](value T, err error) T {
	if err != nil {
		ErrorLog(err.Error())
		return value
	}
	return value
}
func RetrieveError[T any](_ T, err error) error {
	if err != nil {
		ErrorLog(err.Error())
		return err
	}
	return nil
}

/* Log then Panic */
func Panic(message string) {
	message = RetrieveLatestCaller(message)

	fileWriters := []io.Writer{os.Stdout}
	if TraceLogFile != nil {
		fileWriters = append(fileWriters, TraceLogFile)
	}
	if ErrorLogFile != nil {
		fileWriters = append(fileWriters, ErrorLogFile)
	}
	multiWriter := io.MultiWriter(fileWriters...)

	if len(fileWriters) > 0 {
		stackTrace := string(debug.Stack())
		message = stackTrace + "\n" + message
		if _, failed := ErrorExists(fmt.Fprintln(multiWriter, message)); failed {
			fmt.Println("Error writing to the error log during panic, despite the a multiwriter being available")
			return
		}
	} else {
		fmt.Println(message)
	}
	if ErrorLogFile != nil {
		PrintErr(ErrorLogFile.Close())
	}
	if ChangeLogFile != nil {
		PrintErr(ChangeLogFile.Close())
	}
	if TraceLogFile != nil {
		PrintErr(TraceLogFile.Close())
	}
	os.Exit(512)
}
func Assert(value1 any, value2 any) {
	if reflect.TypeOf(value1) != reflect.TypeOf(value2) {
		Panic("Assertion Failed - Types Mismatch `" + reflect.TypeOf(value1).Name() + "` `" + reflect.TypeOf(value2).Name() + "`")
	} else if !reflect.DeepEqual(value1, value2) {
		Panic("Assertion Failed - Values Mismatch `" + fmt.Sprintf("%v", value1) + "` `" + fmt.Sprintf("%v", value2) + "`")
	}
}
func WrapErr(fn func() error) {
	if ErrExists(fn()) {
		ErrorLog("Error During Defer, Continuing...")
	}
}
func WrapPanic(fn func() error) {
	if ErrExists(fn()) {
		Panic("Error During Defer, Exiting...")
	}
}
