# 🪟 Windows Utility Library

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.22-blue?logo=go)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-Windows-lightgrey?logo=windows)]()
[![License](https://img.shields.io/badge/License-MIT-green.svg)]()
[![Status](https://img.shields.io/badge/Status-Experimental-orange)]()

> This library is **Windows-only**. It’s heavily tailored to my personal needs — but if you’re brave enough to use it, here’s what you should know before giving up after the first compiler error.

---

This library provides helper functions similar to what you’d find in other standard libraries, plus some trimmed down standard functions.

It’s intended for almost all projects, however the main purpose is to significantly improve logging and error handling. Yes another one of those.
All error handling writes to the a logfile and stdout.

---

Right after entering `main()`, declare your required variables and call `InitLogs()`:

```go
package main

func main() {
    AppName = "YourApp"
    ExecutableName = "yourexe"
    InitLogs()
}
```

---

You can create a new config with any key->value pair of type: map[string]string. 
```go
WriteConfig(path, yourmap)
```
This will overwrite existing values, but will also leave anything keys that are not in the map you just wrote to file. 
Read the map[string]string config with:
```go
keyvaluepair := ReadConfig()
```

---
Error Handling The main purpose of this library is to smooth out error handling. 
See logging/errors.go for the functions.
Call the error handlers for functions that only return error 
```go
if `ErrExists(thisFunction) {
    // Do something
}
```
This will log the error to error log file and stdout, as well as acting as an operand. 
For functions that return (type, error). 
```go
if var, err := ErrorExists(thisFunction) {
    // Do Something
}
```
Again, this will log the error to both the logfile and stdout. 
The other functions perform similar tasks. Functions with Error expects a return of (T, error) and Err only expects only (error) to be returned. PanicErr(functionMustSucceed); ErrExists(functionMightErrorButWeDontCare); PrintErr(functionMightErrorButOnlyLogIt); All of these functions write to the error log and stdout, Panic() can be called directly to write to an errorlog and os.Exit with a random number.
