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

// CloseLogs writes a final trace entry, disables file logging, and closes all log files.
// Must be used instead of manually closing log files to avoid circular logging on close errors.
func CloseLogs() {
	FileLogging = false
	ConsoleLogging = true

	if ChangeLogFile != nil {
		if _, err := os.Stat(ChangeLogFile.Name()); err != nil {
			ErrorLog("Failed to stat the change log file: " + err.Error())
			ChangeLogFile = nil
			changeLogWriter = nil
		} else {
			if err := ChangeLogFile.Close(); err != nil {
				ErrorLog("Failed to close the change log file: " + err.Error())
			}
			ChangeLogFile = nil
			changeLogWriter = nil
		}
	}
	if TraceLogFile != nil {
		if _, err := os.Stat(TraceLogFile.Name()); err != nil {
			ErrorLog("Failed to stat the trace log file: " + err.Error())
			TraceLogFile = nil
			traceLogWriter = nil
		} else {
			if err := TraceLogFile.Close(); err != nil {
				ErrorLog("Failed to close the trace log file: " + err.Error())
			}
			TraceLogFile = nil
			traceLogWriter = nil
		}
	}
	if ErrorLogFile != nil {
		if _, err := os.Stat(ErrorLogFile.Name()); err != nil {
			ErrorLog("Failed to stat the error log file: " + err.Error())
			ErrorLogFile = nil
			errorLogWriter = nil
		} else {
			if err := ErrorLogFile.Close(); err != nil {
				ErrorLog("Failed to close the error log file: " + err.Error())
			}
			ErrorLogFile = nil
			errorLogWriter = nil
		}
	}
}

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
// - logging.FileLogging (bool) (default: true) - Turns off all file logging, console logging will work if below is unset.
//
// - logging.ConsoleLogging (bool) (default: true)
//
// - logging.RotationCheckInterval (int64 - minutes) (default: 360)
//
// - logging.TraceLogRotation (int64 - days) (default: 3)
//
// - logging.PriorityLogRotation (int64 - days) (default: 14)
func InitLogs() {
	initVars()

	if FileLogging {
		if _, err := os.Stat(BaseLogsFolder); os.IsNotExist(err) {
			PanicErr(os.MkdirAll(BaseLogsFolder, 0755))
		}

		go func() {
			for {
				RotateLogs(BaseLogsFolder)
				time.Sleep(time.Duration(RotationCheckInterval) * time.Minute)
			}
		}()
	}
}
func initVars() {
	if runtime.GOOS == "windows" {
		UserProfile = os.Getenv("USERPROFILE")
		if UserProfile == "" {
			ErrorLog("Unable to locate the user's profile name.")
		}
	} else {
		UserProfile = PanicError(os.UserHomeDir())
		if UserProfile == "" {
			ErrorLog("Unable to locate the user's home directory.")
		}
	}

	if !FileLogging {
		return
	}

	if LoggingPath != "" {
		if runtime.GOOS == "windows" {
			BaseLogsFolder = LoggingPath + "\\"
		} else {
			BaseLogsFolder = LoggingPath + "/"
		}
	} else {
		if runtime.GOOS == "windows" {
			BaseLogsFolder = "C:\\Local\\Logs\\" + AppName + "\\"
		} else {
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
	if !FileLogging {
		return
	}

	currentday := GetDay()

	var ErrorPath string
	var ChangePath string
	var TracePath string

	if foldername == "" {
		ErrorPath = BaseLogsFolder + "errors-" + currentday + ".log"
		ChangePath = BaseLogsFolder + "changes-" + currentday + ".log"
		TracePath = BaseLogsFolder + "trace-" + currentday + ".log"
	} else {
		ErrorPath = BaseLogsFolder + foldername + string(filepath.Separator) + "errors-" + currentday + ".log"
		ChangePath = BaseLogsFolder + foldername + string(filepath.Separator) + "changes-" + currentday + ".log"
		TracePath = BaseLogsFolder + foldername + string(filepath.Separator) + "trace-" + currentday + ".log"
	}
	if _, err := os.Stat(BaseLogsFolder + foldername); os.IsNotExist(err) {
		PanicErr(os.MkdirAll(BaseLogsFolder+foldername, 0755))
	}

	if ErrorLogFile != nil && ErrorPath != ErrorLogFile.Name() {
		PanicErr(ErrorLogFile.Close())
		ErrorLogFile = nil // prevent double close in InitErrorLog
		InitErrorLog(ErrorPath)
	} else if ErrorLogFile == nil {
		InitErrorLog(ErrorPath)
	}
	if ChangeLogFile != nil && ChangePath != ChangeLogFile.Name() {
		PanicErr(ChangeLogFile.Close())
		ChangeLogFile = nil // prevent double close in InitChangeLog
		InitChangeLog(ChangePath)
	} else if ChangeLogFile == nil {
		InitChangeLog(ChangePath)
	}
	if TraceLogFile != nil && TracePath != TraceLogFile.Name() {
		PanicErr(TraceLogFile.Close())
		TraceLogFile = nil // prevent double close in InitTraceLog
		InitTraceLog(TracePath)
	} else if TraceLogFile == nil {
		InitTraceLog(TracePath)
	}
}

// Write to the error log and stdout if ConsoleLogging is set to true.
//
// Example:
//
//	logging.ErrorLog("failed to connect to db")
func ErrorLog(message string) {
	if !FileLogging {
		if ConsoleLogging {
			fmt.Println(message)
		}
		return
	}
	if ErrorLogFile == nil {
		if BaseLogsFolder == "" {
			fmt.Println("BaseLogsFolder is empty, falling back to console-only logging")
			ConsoleLogging = true
			FileLogging = false
			fmt.Println(message)
			return
		}
		currentday := GetDay()
		f, err := os.OpenFile(BaseLogsFolder+"errors-"+currentday+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Failed to open error log file: " + err.Error())
			ConsoleLogging = true
			FileLogging = false
			fmt.Println(message)
			return
		}
		ErrorLogFile = f
		if ConsoleLogging {
			errorLogWriter = io.MultiWriter(os.Stdout, ErrorLogFile)
		} else {
			errorLogWriter = ErrorLogFile
		}
	}

	message = RetrieveLatestCaller(message)
	PrintLogs(message, errorLogWriter)
}

// Write to the change log and stdout if ConsoleLogging is set to true.
//
// Example:
//
//	logging.ChangeLog("updated user profile", "user-123")
func ChangeLog(message string, idnumber string) {
	if !FileLogging {
		if ConsoleLogging {
			fmt.Println(message)
		}
		return
	}
	if ChangeLogFile == nil {
		if BaseLogsFolder == "" {
			fmt.Println("BaseLogsFolder is empty, falling back to console-only logging")
			ConsoleLogging = true
			FileLogging = false
			fmt.Println(message)
			return
		}
		currentday := GetDay()
		f, err := os.OpenFile(BaseLogsFolder+"changes-"+currentday+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Failed to open change log file: " + err.Error())
			ConsoleLogging = true
			FileLogging = false
			fmt.Println(message)
			return
		}
		ChangeLogFile = f
		if ConsoleLogging {
			changeLogWriter = io.MultiWriter(os.Stdout, ChangeLogFile)
		} else {
			changeLogWriter = ChangeLogFile
		}
	}

	message = RetrieveLatestCaller(message)
	if idnumber != "" {
		message = message + " || " + idnumber
	}
	PrintLogs(message, changeLogWriter)
}

// Write to the trace log and stdout if ConsoleLogging is set to true.
//
// Example:
//
//	logging.TraceLog("starting background job")
func TraceLog(message string) {
	if !FileLogging {
		if ConsoleLogging {
			fmt.Println(message)
		}
		return
	}
	if !TraceDebug {
		return
	}
	if TraceLogFile == nil {
		if BaseLogsFolder == "" {
			fmt.Println("BaseLogsFolder is empty, falling back to console-only logging")
			ConsoleLogging = true
			FileLogging = false
			fmt.Println(message)
			return
		}
		currentday := GetDay()
		f, err := os.OpenFile(BaseLogsFolder+"trace-"+currentday+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Failed to open trace log file: " + err.Error())
			ConsoleLogging = true
			FileLogging = false
			fmt.Println(message)
			return
		}
		TraceLogFile = f
		if ConsoleLogging {
			traceLogWriter = io.MultiWriter(os.Stdout, TraceLogFile)
		} else {
			traceLogWriter = TraceLogFile
		}
	}
	message = RetrieveLatestCaller(message)
	PrintLogs(message, traceLogWriter)
}

// RetrieveLatestCaller formats a log message with the first caller outside the logging package.
//
// Example:
//
//	withCaller := logging.RetrieveLatestCaller("starting worker")
func RetrieveLatestCaller(message string) string {
	for skip := 1; skip <= 10; skip++ {
		pc, _, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fullName := runtime.FuncForPC(pc).Name()

		// Extract the short package.Function portion after the last /
		shortName := fullName
		if lastSlash := strings.LastIndex(fullName, "/"); lastSlash >= 0 {
			shortName = fullName[lastSlash+1:]
		}
		if strings.HasPrefix(shortName, "logging.") {
			continue
		}

		// Split on the last dot to separate package path from function name
		lastDot := strings.LastIndex(fullName, ".")
		pkg := fullName[:lastDot]
		funcName := fullName[lastDot+1:]
		if idx := strings.Index(funcName, "["); idx >= 0 {
			funcName = funcName[:idx]
		}

		return fmt.Sprintf("%s || %-60s || %s", GetTime(), fmt.Sprintf("(%s) %s:%d", pkg, funcName, line), message)
	}
	return fmt.Sprintf("%s || NO CALLER || %s", GetTime(), message)
}

// PrintLogs writes the message to the requested log stream.
//
// Example:
//
//	logging.PrintLogs("system online", traceLogWriter)
func PrintLogs(message string, writer io.Writer) {
	if _, err := fmt.Fprintln(writer, message); err != nil {
		fmt.Println("File logging failed due to an error, reverting to console logging.")
		fmt.Println(message)
		ConsoleLogging = true
		FileLogging = false
		return
	}
}

// Rotate the logs in the folder, removes logs older than the rotation time set in TraceLogRotation and PriorityLogRotation
//
// Example Usage:
//
//	RotateLogs(BaseLogsFolder) // Immediately rotate - Change the folder path for subfolder rotation
func RotateLogs(logFolder string) {
	if !FileLogging {
		return
	}
	files, err := os.ReadDir(logFolder)
	if err != nil || len(files) == 0 {
		return
	} else {
		for _, file := range files {
			if file.IsDir() {
				RotateLogs(filepath.Join(logFolder, file.Name()))
				continue
			}

			logSplit := strings.Split(file.Name(), "-")
			if len(logSplit) < 4 {
				continue
			}
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
	if !FileLogging {
		return
	}
	if (ErrorLogFile != nil && fullPath == ErrorLogFile.Name()) ||
		(ChangeLogFile != nil && fullPath == ChangeLogFile.Name()) ||
		(TraceLogFile != nil && fullPath == TraceLogFile.Name()) {
		return
	}
	TraceLog("Removing log file: " + fullPath)
	PrintErr(os.Remove(fullPath))
}
