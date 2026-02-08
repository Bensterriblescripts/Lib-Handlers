package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	. "github.com/Bensterriblescripts/Lib-Handlers/time"
)

var errorLogWriter io.Writer
var changeLogWriter io.Writer
var traceLogWriter io.Writer

// Initialize all logging files and variables
//
// # The following must be set before calling:
//
// - logging.AppName (string)
//
// - logging.LoggingPath (string) (default: C:\Local\Logs\AppName\ or /home/*user*/local/Logs/AppName/)
//
// - logging.TraceDebug (bool) (default: true)
//
// - logging.ConsoleLogging (bool) (default: true)
//
// - logging.TraceLogRotation (int64 - days) (default: 3)
//
// - logging.PriorityLogRotation (int64 - days) (default: 14)
func InitLogs() {
	initVars()

	if _, err := os.Stat(BaseLogsFolder); os.IsNotExist(err) {
		PanicErr(os.MkdirAll(BaseLogsFolder, 0755))
	}
	go func() {
		for {
			RotateLogs(BaseLogsFolder)
			time.Sleep(time.Duration(TraceLogRotation) * time.Minute)
		}
	}()
}
func initVars() {
	if LoggingPath != "" {
		if runtime.GOOS == "windows" {
			BaseLogsFolder = LoggingPath + "\\"
		} else {
			BaseLogsFolder = LoggingPath + "/"
		}
	} else {
		if runtime.GOOS == "windows" {
			UserProfile = os.Getenv("USERPROFILE")
			if UserProfile == "" {
				Panic("Unable to locate the user's profile name.")
			}
			BaseLogsFolder = "C:\\Local\\Logs\\" + AppName + "\\"
		} else {
			UserProfile = PanicError(os.UserHomeDir())
			BaseLogsFolder = UserProfile + "/local/Logs/" + AppName + "/"
		}
	}
}

// Initialize the error log file, the file will be closed and reopened if it already exists.
//
// Example:
//
//	logging.InitErrorLog("C:\\Local\\Logs\\app\\errors.log")
//	defer logging.ErrorLogFile.Close()
//
// Sets ErrorLogFile
func InitErrorLog(filename string) {
	if ErrorLogFile == nil {
		ErrorLogFile = PanicError(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	} else {
		PrintErr(ErrorLogFile.Close())
		ErrorLogFile = PanicError(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	}
	if ConsoleLogging {
		errorLogWriter = io.MultiWriter(os.Stdout, ErrorLogFile)
	} else {
		errorLogWriter = ErrorLogFile
	}
}

// Initialize the change log file, the file will be closed and reopened if it already exists.
//
// Example:
//
//	logging.InitChangeLog("C:\\Local\\Logs\\app\\changes.log")
//	defer logging.ChangeLogFile.Close()
//
// Sets ChangeLogFile
func InitChangeLog(filename string) {
	if ChangeLogFile == nil {
		ChangeLogFile = PanicError(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	} else {
		PrintErr(ChangeLogFile.Close())
		ChangeLogFile = PanicError(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	}

	if ConsoleLogging {
		changeLogWriter = io.MultiWriter(os.Stdout, ChangeLogFile)
	} else {
		changeLogWriter = ChangeLogFile
	}
}

// Initialize the trace log file, the trace log will not be written to if TraceDebug is manually set to false.
//
// Example:
//
//	logging.InitTraceLog("C:\\Local\\Logs\\app\\trace.log")
//	defer logging.TraceLogFile.Close()
//
// Sets TraceLogFile
func InitTraceLog(filename string) {
	if !TraceDebug {
		return
	}
	if TraceLogFile == nil {
		TraceLogFile = PanicError(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	} else {
		PrintErr(TraceLogFile.Close())
		TraceLogFile = PanicError(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	}

	if ConsoleLogging {
		traceLogWriter = io.MultiWriter(os.Stdout, TraceLogFile)
	} else {
		traceLogWriter = TraceLogFile
	}
}

// Alter the current log folder path.
//
// Example:
//
//	logging.SetLogsFolder("session-2024-01-01")
//
// The new logs folder will sit under AppName/StringPassedIn
func SetLogsFolder(foldername string) {
	currentday := GetDay()

	var ErrorPath string
	var ChangePath string
	var TracePath string

	if foldername == "" {
		ErrorPath = BaseLogsFolder + "errors-" + currentday + ".log"
		ChangePath = BaseLogsFolder + "changes-" + currentday + ".log"
		TracePath = BaseLogsFolder + "trace-" + currentday + ".log"
	} else {
		ErrorPath = BaseLogsFolder + foldername + "\\" + "errors-" + currentday + ".log"
		ChangePath = BaseLogsFolder + foldername + "\\" + "changes-" + currentday + ".log"
		TracePath = BaseLogsFolder + foldername + "\\" + "trace-" + currentday + ".log"
	}
	if _, err := os.Stat(BaseLogsFolder + foldername); os.IsNotExist(err) {
		PanicErr(os.MkdirAll(BaseLogsFolder+foldername, 0755))
	}

	if ErrorLogFile != nil && ErrorPath != ErrorLogFile.Name() {
		PanicErr(ErrorLogFile.Close())
		InitErrorLog(ErrorPath)
	} else if ErrorLogFile == nil {
		InitErrorLog(ErrorPath)
	}
	if ChangeLogFile != nil && ChangePath != ChangeLogFile.Name() {
		PanicErr(ChangeLogFile.Close())
		InitChangeLog(ChangePath)
	} else if ChangeLogFile == nil {
		InitChangeLog(ChangePath)
	}
	if TraceLogFile != nil && TracePath != TraceLogFile.Name() {
		PanicErr(TraceLogFile.Close())
		InitTraceLog(TracePath)
	}
}

// Write to the error log and stdout if ConsoleLogging is set to true.
//
// Example:
//
//	logging.ErrorLog("failed to connect to db")
func ErrorLog(message string) {
	if ErrorLogFile == nil {
		currentday := GetDay()
		InitErrorLog(BaseLogsFolder + "errors-" + currentday + ".log")
		if ErrorLogFile == nil {
			Panic("ErrorLogFile is nil")
		}
	}

	message = RetrieveLatestCaller(message)
	PrintLogs(message, 0)
}

// Write to the change log and stdout if ConsoleLogging is set to true.
//
// Example:
//
//	logging.ChangeLog("updated user profile", "user-123")
func ChangeLog(message string, idnumber string) {
	if ChangeLogFile == nil {
		currentday := GetDay()
		InitChangeLog(BaseLogsFolder + "changes-" + currentday + ".log")
		if ChangeLogFile == nil {
			Panic("ChangeLogFile is nil")
		}
	}

	message = RetrieveLatestCaller(message)
	if idnumber != "" {
		message = message + " || " + idnumber
	}
	PrintLogs(message, 1)
}

// Write to the trace log and stdout if ConsoleLogging is set to true.
//
// Example:
//
//	logging.TraceLog("starting background job")
func TraceLog(message string) {
	if !TraceDebug {
		return
	}
	if TraceLogFile == nil {
		currentday := GetDay()
		InitTraceLog(BaseLogsFolder + "trace-" + currentday + ".log")
		if TraceLogFile == nil {
			Panic("TraceLogFile is nil")
		}
	}
	message = RetrieveLatestCaller(message)
	PrintLogs(message, 2)
}

// RetrieveLatestCaller formats a log message with caller details.
//
// Example:
//
//	withCaller := logging.RetrieveLatestCaller("starting worker")
func RetrieveLatestCaller(message string) string {
	pc, _, callerline, ok := runtime.Caller(1)
	if !ok {
		return fmt.Sprintf("%s || NO CALLER || %s", GetTime(), message)
	}
	caller := strings.Split(runtime.FuncForPC(pc).Name(), ".")

	pc2, _, callerline2, ok2 := runtime.Caller(2)
	if !ok2 {
		return fmt.Sprintf("%s || %-60s || %s", GetTime(), fmt.Sprintf("(%s) %s:%d", caller[0], caller[1], callerline), message)
	}
	caller2 := strings.Split(runtime.FuncForPC(pc2).Name(), ".")
	caller2[1] = strings.Replace(caller2[1], "[", "", 1)

	pc3, _, callerline3, ok3 := runtime.Caller(3)
	if !ok3 {
		return fmt.Sprintf("%s || %-60s || %s", GetTime(), fmt.Sprintf("(%s) %s:%d", caller2[0], caller2[1], callerline2), message)
	}
	caller3 := strings.Split(runtime.FuncForPC(pc3).Name(), ".")
	caller3[1] = strings.Replace(caller3[1], "[", "", 1)

	return fmt.Sprintf("%s || %-60s || %s", GetTime(), fmt.Sprintf("(%s) %s:%d", caller3[0], caller3[1], callerline3), message)
}

// PrintLogs writes the message to the requested log stream.
//
// Example:
//
//	logging.PrintLogs("system online", 2)
func PrintLogs(message string, errorlevel int) {
	switch errorlevel {
	case 0:
		PanicError(fmt.Fprintln(errorLogWriter, message))
	case 1:
		PanicError(fmt.Fprintln(changeLogWriter, message))
	case 2:
		PanicError(fmt.Fprintln(traceLogWriter, message))
	}
}

// Rotate the logs in the folder, removes logs older than the rotation time set in TraceLogRotation and PriorityLogRotation
//
// Example Usage:
//
//	RotateLogs(BaseLogsFolder) // Immediately rotate - Change the folder path for subfolder rotation
func RotateLogs(logFolder string) {
	files, err := os.ReadDir(logFolder)
	if err != nil || len(files) == 0 {
		return
	} else {
		for _, file := range files {
			if file.IsDir() {
				RotateLogs(logFolder + file.Name())
				continue
			}

			logSplit := strings.Split(file.Name(), "-")
			logType := logSplit[0]
			logYear := logSplit[1]
			logMonth := logSplit[2]
			logDay := logSplit[3]
			fullPath := filepath.Join(logFolder, file.Name())

			switch logType {
			case "trace":
				clearOutdatedLogs(fullPath, logYear, logMonth, logDay, int(TraceLogRotation))
			case "changes", "errors":
				clearOutdatedLogs(fullPath, logYear, logMonth, logDay, int(PriorityLogRotation))
			default:
				ErrorLog("Unknown log type: " + logType)
			}
		}
	}
}

// clearOutdatedLogs removes logs older than the provided retention window.
//
// Example:
//
//	clearOutdatedLogs("C:\\Local\\Logs\\app\\trace-2024-1-1.log", "2024", "1", "1", 3)
func clearOutdatedLogs(fullPath string, logStringYear string, logStringMonth string, logStringDay string, daysToKeep int) {
	if logStringYear == "" {
		ErrorLog("Log file error, year is empty: " + fullPath)
		return
	}
	if logStringMonth == "" {
		ErrorLog("Log file error, month is empty: " + fullPath)
		return
	}
	if logStringDay == "" {
		ErrorLog("Log file error, day is empty: " + fullPath)
		return
	}
	logStringDay = strings.Replace(logStringDay, ".log", "", 1)
	logDay, errDay := ErrorExists(strconv.Atoi(logStringDay))
	if errDay {
		ErrorLog("Log file error, day is not a number: " + fullPath)
		return
	}
	logMonth, errMonth := ErrorExists(strconv.Atoi(logStringMonth))
	if errMonth {
		ErrorLog("Log file error, month is not a number: " + fullPath)
		return
	}
	logYear, errYear := ErrorExists(strconv.Atoi(logStringYear))
	if errYear {
		ErrorLog("Log file error, year is not a number: " + fullPath)
		return
	}

	logDate := time.Date(logYear, time.Month(logMonth), logDay, 0, 0, 0, 0, time.UTC)
	cutoff := time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 0, -daysToKeep)

	if logDate.Before(cutoff) {
		RemoveLog(fullPath)
	}
}

// RemoveLog deletes the log file at the provided path.
//
// Example:
//
//	logging.RemoveLog("C:\\Local\\Logs\\app\\trace-2024-1-1.log")
func RemoveLog(fullPath string) {
	TraceLog("Removing log file: " + fullPath)
	PrintErr(os.Remove(fullPath))
}
