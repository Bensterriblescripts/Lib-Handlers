package logging

import (
	"os"
)

/* Set these ones in main */
var AppName string = "Default"
var TraceDebug bool = true
var FileLogging bool = true
var ConsoleLogging bool = true
var LoggingPath string = ""

/* Automatically Set and Optional */
var UserProfile string
var RotationCheckInterval int64 = 360 // Time between rotation checks - minutes
var TraceLogRotation int64 = 3        // Max age of trace logs - days
var PriorityLogRotation int64 = 14    // Max age of change/error logs - days

var BaseLogsFolder string
var ErrorLogFile *os.File
var ChangeLogFile *os.File
var TraceLogFile *os.File
