package logging

import (
	"os"
)

/* Set these ones in main */
var AppName string = "Default"
var ExecutableName string = "Default"
var TraceDebug bool = false
var NetworkDebug bool = false
var FileLogging bool = true

/* Optional */
var OPENAIApiKey string

/* Automatically Set and Optional */
var UserProfile string
var ConfigPath string = "C:\\Local\\Config\\Default.ini"
var CacheFrequency int64 = 12 // Minutes between cache refreshes
var ProgramLogString string
var TraceLogRotation int64 = 3     // Time between trace log rotation - days
var PriorityLogRotation int64 = 14 // Time between change/error log rotation - days

var BaseLogsFolder string
var ErrorLogFile *os.File
var ChangeLogFile *os.File
var TraceLogFile *os.File
