For any poor sod who decides to use this library, I tailored it to my own needs, here is a short message so you don't give up after the first compiler error.
Heavy warning, this library is for Windows only.

It contains helper functions that you might find in other language standard libraries or reduce the error handling required for existing functions.


/** Right after entering main() **/
Declare InitLogs() within your main entry point after declaring the variables noted in logging/vars.go.
As of writing, the variables required before InitLogs() are:
var AppName string = "Default"
var ExecutableName string = "Default"
var TraceDebug bool = false
var NetworkDebug bool = false
var FileLogging bool = true
You MUST set appname and executable name.

All logs will then be directed to C:\Local\*AppName*\(trace|change|error).log
All configuration writes are stored in in C:\Local\Config\*AppName*.ini 

You can create a new config with any key->value pair of type: map[string]string.
WriteConfig(path, yourmap)
This will overwrite existing values, but will also leave anything keys that are not in the map you just wrote to file.

Read the map[string]string config with:
keyvaluepair := ReadConfig()

/** Error Handling **/
The main purpose of this library is to smooth out error handling.
See logging/errors.go for the functions

Call the error handlers like so for functions that only return error
if ErrExists(thisFunction) {
  // Do something
}
This will log the error to error log file and stdout, as well as acting as an operand.

For functions that return (type, error).
if var, err := ErrorExists(thisFunction) {
  // Do Something
}
Again, this will log the error to both the logfile and stdout.

The other functions perform similar tasks. Functions with Error expects a return of (T, error) and Err only expects only (error) to be returned.

PanicErr(functionMustSucceed);
ErrExists(functionMightErrorButWeDontCare);
PrintErr(functionMightErrorButOnlyLogIt);

All of these functions write to the error log and stdout, Panic() can be called directly to write to an errorlog and os.Exit with a random number.
