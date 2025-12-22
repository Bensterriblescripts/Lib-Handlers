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
			// time.Sleep(30 * time.Second)
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

// Initialize the error log file, the file will be closed and reopened if it already exists
//
// Sets ErrorLogFile
func InitErrorLog(filename string) {
	if ErrorLogFile == nil {
		ErrorLogFile = PanicError(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	} else {
		ErrorLogFile.Close()
		ErrorLogFile = PanicError(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	}
	if ConsoleLogging {
		errorLogWriter = io.MultiWriter(os.Stdout, ErrorLogFile)
	} else {
		errorLogWriter = ErrorLogFile
	}
}

// Initialize the change log file, the file will be closed and reopened if it already exists
//
// Sets ChangeLogFile
func InitChangeLog(filename string) {
	if ChangeLogFile == nil {
		ChangeLogFile = PanicError(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	} else {
		ChangeLogFile.Close()
		ChangeLogFile = PanicError(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	}

	if ConsoleLogging {
		changeLogWriter = io.MultiWriter(os.Stdout, ChangeLogFile)
	} else {
		changeLogWriter = ChangeLogFile
	}
}

// Initialize the trace log file, the trace log will not be written to if TraceDebug is manually set to false
//
// Sets TraceLogFile
func InitTraceLog(filename string) {
	if !TraceDebug {
		return
	}
	if TraceLogFile == nil {
		TraceLogFile = PanicError(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	} else {
		TraceLogFile.Close()
		TraceLogFile = PanicError(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	}

	if ConsoleLogging {
		traceLogWriter = io.MultiWriter(os.Stdout, TraceLogFile)
	} else {
		traceLogWriter = TraceLogFile
	}
}

// Alter the current log folder path
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

// Write to the error log and stdout if ConsoleLogging is set to true
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

// Write to the change log and stdout if ConsoleLogging is set to true
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

// Write to the trace log and stdout if ConsoleLogging is set to true
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
//	go func() {
//		for {
//			RotateLogs(BaseLogsFolder) // Immediately rotate - Change the folder path for subfolder rotation
//			time.Sleep(time.Duration(TraceLogRotation) * time.Day) // Wait for the next rotation in long running apps
//		}
//	}()
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
				ClearOutdatedLogs(fullPath, logYear, logMonth, logDay, int(TraceLogRotation))
			case "changes", "errors":
				ClearOutdatedLogs(fullPath, logYear, logMonth, logDay, int(PriorityLogRotation))
			default:
				ErrorLog("Unknown log type: " + logType)
			}
		}
	}
}
func ClearOutdatedLogs(fullPath string, logStringYear string, logStringMonth string, logStringDay string, daysToKeep int) {
	if logStringYear == "" {
		ErrorLog("Log file error, year is empty: " + fullPath)
	}
	if logStringMonth == "" {
		ErrorLog("Log file error, month is empty: " + fullPath)
	}
	if logStringDay == "" {
		ErrorLog("Log file error, day is empty: " + fullPath)
	}
	logDayClean := strings.ReplaceAll(logStringDay, ".log", "")
	logDay, err := ErrorExists(strconv.Atoi(logDayClean))
	if err {
		ErrorLog("Log file error, day is not a number: " + fullPath)
	}
	logMonth, err := ErrorExists(strconv.Atoi(logStringMonth))
	if err {
		ErrorLog("Log file error, month is not a number: " + fullPath)
	}
	logYear, err := ErrorExists(strconv.Atoi(logStringYear))
	if err {
		ErrorLog("Log file error, year is not a number: " + fullPath)
	}

	currentDateArray := GetDateArray()
	currentDay := currentDateArray[0]
	currentMonth := currentDateArray[1]
	currentYear := currentDateArray[2]

	if logDay > daysToKeep { // If our days don't go back to the previous month
		if logYear != currentYear || logMonth != currentMonth {
			RemoveLog(fullPath)
		}
		if logDay < currentDay {
			RemoveLog(fullPath)
		}
	}
}
func RemoveLog(fullPath string) {
	PrintErr(os.Remove(fullPath))
}
