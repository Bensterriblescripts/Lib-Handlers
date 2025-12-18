package logging

import (
	"os"
)

/* Set these ones in main */
var AppName string = "Default"
var TraceDebug bool = true
var ConsoleLogging bool = true
var LoggingPath string = ""

/* Automatically Set and Optional */
var UserProfile string
var TraceLogRotation int64 = 3     // Time between trace log rotation - days
var PriorityLogRotation int64 = 14 // Time between change/error log rotation - days

var BaseLogsFolder string
var ErrorLogFile *os.File
var ChangeLogFile *os.File
var TraceLogFile *os.File
