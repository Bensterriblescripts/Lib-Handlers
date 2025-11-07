# 🪟 Windows Utility Library

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.22-blue?logo=go)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-Windows-lightgrey?logo=windows)]()
[![License](https://img.shields.io/badge/License-MIT-green.svg)]()
[![Status](https://img.shields.io/badge/Status-Experimental-orange)]()

> This library is **Windows-only**. It’s heavily tailored to my personal needs — but if you’re brave enough to use it, here’s what you should know before giving up after the first compiler error.

---

This library is focused is on error handling with logs, but also has functions I use regularly.
Such as:
 - OAuth2 access token retrieval
 - Dataverse odata requests
 - Guid types
 - SSH Tunneling and key access
 - Windows API HWND retrieval and changes (for manipulating active program's window state)


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
If there is a previous config, it will overwrite the existing values and retain the values that were not present in the map parameter. 
Read the map[string]string config with:
```go
keyvaluepair := ReadConfig()
```

---
Error Handling The main purpose of this library is to smooth out error handling. 
See logging/errors.go for the functions.
Call the error handlers for functions that only return error 
```go
if ErrExists(thisFunction) {
    // Handle error, log has already been created.
}
```
This will log the error to error log file and stdout, as well as acting as an operand. 
For functions that return (type, error). 
```go
if var, failed := ErrorExists(thisFunction); failed {
    // Handle error, log has already been created.
} else {
    // Use our var as normal
}
```
Again, this will log the error to both the logfile and stdout. 
The other functions perform similar tasks.

Functions with Error expects a return of (T, error) and Err only expects only (error) to be returned. 
`PanicErr(functionMustSucceed);` `ErrExists(functionMightErrorButWeDontCare);` `PrintErr(functionMightErrorButOnlyLogIt);` 

---

All of these functions write to the error log and stdout, `Panic()` can be called directly to write to an errorlog and exit the program.
YOU MUST HANDLE YOUR DEFERS before using this function.
