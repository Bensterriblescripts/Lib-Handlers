# 🪟 Windows Utility Library for Golang

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.25.3-blue?logo=go)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-Windows-lightgrey?logo=windows)]()
[![License](https://img.shields.io/badge/License-MIT-green.svg)]()
[![Status](https://img.shields.io/badge/Status-Experimental-orange)]()

> This library is **Windows-only** and heavily tailored to my personal needs.  I've detailed the main uses and initial setup below just in case this is in use elsewhere.

---

This library is focused is on error handling with logs, but also has functions I use regularly.
Such as:
 - OAuth2 access token retrieval
 - Dataverse odata requests
 - Guid types
 - SSH Tunneling and key access
 - Windows API HWND retrieval and changes (for manipulating active program's window state)
 - OpenAI request format (responses endpoint, not completions)

---

run `go get github.com/Bensterriblescripts/Lib-Handlers@latest`

or

run `go get github.com/Bensterriblescripts/Lib-Handlers@linux`

---

Right after entering `main()`, you need to declare some mandatory variables as well as `InitLogs()`:

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
WriteConfig(yourmap)
```
If there is a previous config, it will overwrite the existing values and retain the values that were not present in the map parameter. 
Read the map[string]string config with:
```go
keyvaluepair := ReadConfig()
```

---
**Error Handling**

This here is main purpose of the library.
Golangs error handling sucks, you cannot tell me otherwise. 

All calls to these functions writes into C:\Local\Logs\*Appname*\ as well as stdout.
*Functions with Err expect only (error), functions with Error expect (T, error)*
```go
if ErrExists(thisFunction) {
    // Handle error, log has already been created.
}
```
This will log the error to C:\Local\Logs\\*AppName\*\error-datestring.log and stdout.

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
 - `PanicErr(functionMustSucceed);`
 - `ErrExists(functionShouldBeCheckedForFailing);`
 - `PrintErr(justLogTheErrorWeDontCare);`
 
 - `PanicError(functionMustSucceed);`
 - `ErrorExists(functionShouldBeCheckedForFailing);`
 - `PrintError(justLogTheErrorWeDontCare);`
 - `defer WrapErr(thingWeWantToDefer)`

---

All of these functions write to the error log and stdout, `Panic()` can be called directly to write to an errorlog and exit the program.
YOU MUST HANDLE YOUR DEFERS before using this function.
