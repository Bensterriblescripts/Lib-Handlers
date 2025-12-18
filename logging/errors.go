package logging

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime/debug"
)

// Panic if an error is returned
//
//	func() {
//		err := functionName() // Example - Function only returns an error
//		PanicErr(functionName())
//	}
func PanicErr(err error) {
	if err != nil {
		Panic(err.Error())
	}
}

// Log the error and continue
//
//	func() {
//		err := functionName() // Example - Function only returns an error
//		PrintErr(functionName())
//	}
func PrintErr(err error) {
	if err != nil {
		ErrorLog(err.Error())
	}
}

// Check if an error exists
//
//	func() {
//		err := functionName() // Example - Function only returns an error
//		if ErrExists(functionName()) {
//			// Handle error, log has already been created.
//		}
//	}
func ErrExists(err error) bool {
	if err != nil {
		ErrorLog(err.Error())
		return true
	}
	return false
}

// Panic if any error is returned
//
//	func() {
//		value, err := functionName() // Example - Function returns a value and an error
//		val := PanicError(functionName())
//	}
func PanicError[T any](value T, err error) T {
	if err != nil {
		Panic(err.Error())
	}
	return value
}

// Log the error and continue
//
//	func() {
//		value, err := functionName() // Example - Function returns a value and an error
//		val := PrintError(functionName())
//	}
func PrintError[T any](value T, err error) T {
	if err != nil {
		ErrorLog(err.Error())
	}
	return value
}

// Check if an error exists and log it
//
//	func() {
//		value, err := functionName() // Example - Function returns a value and an error
//		if val, exists := ErrorExists(functionName()); exists {
//			// Handle error, log has already been created.
//		}
//		else {
//			// Use our value as normal
//		}
//	}
func ErrorExists[T any](value T, err error) (T, bool) {
	if err != nil {
		ErrorLog(err.Error())
		return value, true
	}
	return value, false
}

// Wrap an error in a defer function
//
//	func() {
//		defer WrapErr(functionName)
//	}
func WrapErr(fn func() error) {
	if ErrExists(fn()) {
		ErrorLog("Error During Defer, Continuing...")
	}
}

// Wrap an error in a defer function and panic if an error occurs
//
//	func() {
//		defer WrapPanic(functionName)
//	}
func WrapPanic(fn func() error) {
	if ErrExists(fn()) {
		Panic("Error During Defer, Exiting...")
	}
}

// Log then Panic
//
//	func() {
//		Panic("critical failure")
//	}
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

// Assert that two values are of the same type and value
//
//	func() {
//		Assert(10, 10)
//	}
func Assert(value1 any, value2 any) {
	if reflect.TypeOf(value1) != reflect.TypeOf(value2) {
		Panic("Assertion Failed - Types Mismatch `" + reflect.TypeOf(value1).Name() + "` `" + reflect.TypeOf(value2).Name() + "`")
	} else if !reflect.DeepEqual(value1, value2) {
		Panic("Assertion Failed - Values Mismatch `" + fmt.Sprintf("%v", value1) + "` `" + fmt.Sprintf("%v", value2) + "`")
	}
}
