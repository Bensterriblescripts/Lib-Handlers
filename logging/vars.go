package logging

import (
	"os"
)

var AppName string = "Default"
var ExecutableName string = "Default"
var UserProfile string
var DevelopmentMode bool = false

var CacheFrequency int64 = 12 // Minutes between cache refreshes

var ProgramLogString string
var TraceLogRotation int64 = 3     // Time between trace log rotation - days
var PriorityLogRotation int64 = 14 // Time between change/error log rotation - days

var FileLogging bool = true
var BaseLogsFolder string
var ErrorLogFile *os.File
var ChangeLogFile *os.File
var TraceLogFile *os.File
