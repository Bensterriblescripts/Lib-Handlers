package logging

import (
	"fmt"
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
		ErrorLog("Error incurred during defer, continuing...")
	}
}

// Wrap an error in a defer function and panic if an error occurs
//
//	func() {
//		defer WrapPanic(functionName)
//	}
func WrapPanic(fn func() error) {
	if ErrExists(fn()) {
		Panic("Error incurred during defer, panic exiting...")
	}
}

// Log then Panic
//
//	func() {
//		Panic("critical failure")
//	}
//
// Closes all log files and exits with code 512
func Panic(message string) {
	message = RetrieveLatestCaller(message)
	stackTrace := string(debug.Stack())

	fmt.Println(stackTrace + "\n" + message)

	if FileLogging && ErrorLogFile != nil {
		if _, err := fmt.Fprintln(ErrorLogFile, stackTrace+"\n"+message); err != nil {
			fmt.Println("Error writing to the errorlog during panic.")
		}
	}
	CloseLogs()
	os.Exit(512)
}

// Assert that two values are of the same type and value
//
//	func() {
//		Assert(10, 10)
//	}
func Assert(value1 any, value2 any) {
	if value1 == nil || value2 == nil {
		if value1 == nil && value2 == nil {
			return
		}
		Panic("Assertion Failed - One value is nil, the other is not")
		return
	}
	if reflect.TypeOf(value1) != reflect.TypeOf(value2) {
		Panic("Assertion Failed - Types Mismatch `" + reflect.TypeOf(value1).Name() + "` `" + reflect.TypeOf(value2).Name() + "`")
	} else if !reflect.DeepEqual(value1, value2) {
		Panic("Assertion Failed - Values Mismatch `" + fmt.Sprintf("%v", value1) + "` `" + fmt.Sprintf("%v", value2) + "`")
	}
}
