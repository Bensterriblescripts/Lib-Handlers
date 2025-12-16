package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	. "github.com/Bensterriblescripts/Lib-Handlers/time"
)

var errorLogWriter io.Writer
var changeLogWriter io.Writer
var traceLogWriter io.Writer

func InitVars() {
	if runtime.GOOS == "windows" {
		UserProfile = os.Getenv("USERPROFILE")
		if UserProfile == "" {
			Panic("Unable to locate the user's profile name.")
		}
		BaseLogsFolder = "C:\\Local\\Logs\\" + AppName + "\\"
	} else {
		UserProfile = PanicError(os.UserHomeDir())
		BaseLogsFolder = UserProfile + "/logs/" + AppName + "/"
	}
}
func InitLogs() {
	InitVars()

	if _, err := os.Stat(BaseLogsFolder); os.IsNotExist(err) {
		PanicErr(os.MkdirAll(BaseLogsFolder, 0755))
	}
	go func() {
		for {
			time.Sleep(30 * time.Second)
			RotateLogs(BaseLogsFolder)
			time.Sleep(time.Duration(TraceLogRotation) * time.Minute)
		}
	}()
}
func InitErrorLog(filename string) {
	ErrorLogFile = PanicError(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	if ConsoleLogging {
		errorLogWriter = io.MultiWriter(os.Stdout, ErrorLogFile)
	} else {
		errorLogWriter = ErrorLogFile
	}
}
func InitChangeLog(filename string) {
	ChangeLogFile = PanicError(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	if ConsoleLogging {
		changeLogWriter = io.MultiWriter(os.Stdout, ChangeLogFile)
	} else {
		changeLogWriter = ChangeLogFile
	}
}
func InitTraceLog(filename string) {
	if !TraceDebug {
		return
	}
	TraceLogFile = PanicError(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	if ConsoleLogging {
		traceLogWriter = io.MultiWriter(os.Stdout, TraceLogFile)
	} else {
		traceLogWriter = TraceLogFile
	}
}
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
func RotateLogs(logFolder string) {
	if TraceDebug {
		return
	}
	files, err := os.ReadDir(logFolder)
	if err != nil || len(files) == 0 {
		return
	} else {
		for _, file := range files {
			if file.IsDir() {
				RotateLogs(logFolder + file.Name())
				continue
			}

			logtype := strings.Split(file.Name(), "-")[0]
			cleandate := strings.ReplaceAll(file.Name(), logtype+"-", "")
			cleandate = strings.ReplaceAll(cleandate, ".log", "")

			switch logtype {
			case "trace":
				if LogOutdated(cleandate, int64(TraceLogRotation)) {
					ClearOutdatedLogs(filepath.Join(logFolder, file.Name()))
				}
			case "changes", "errors":
				if LogOutdated(cleandate, int64(PriorityLogRotation)) {
					ClearOutdatedLogs(filepath.Join(logFolder, file.Name()))
				}
			default:
				ErrorLog("Unknown log type: " + logtype)
			}
		}
	}
}
func LogOutdated(filedate string, daysToKeep int64) bool {
	if filedateint, err := ErrorExists(time.Parse("2006-1-2", filedate)); !err {
		filedateunix := filedateint.Unix()
		return filedateunix < GetUnixTime()-daysToKeep
	}
	return true
}
func ClearOutdatedLogs(file string) {
	TraceLog("Clearing outdated log file: " + file)
	PrintErr(os.Remove(file))
}
