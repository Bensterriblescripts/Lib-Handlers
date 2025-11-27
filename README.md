# 🪟 Windows Utility Library for Golang

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.25.3-blue?logo=go)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-Windows-lightgrey?logo=windows)]()
[![License](https://img.shields.io/badge/License-MIT-green.svg)]()
[![Status](https://img.shields.io/badge/Status-Experimental-orange)]()

> This library is **Windows-only** and heavily tailored to my personal needs.  I've detailed the main uses and initial setup below just in case this is in use elsewhere.
> There is a Linux branch which trims out functions which use Windows DLLs and API calls.
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

BE WARNED: These functions are innefficient and each call contains a deep copy. Avoid using these in production where possible.

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
You must add a defer to close the files and any ssh handlers before the program closes. I am too stupid/lazy to figure out how to do this outside of the main function at this moment.

---

## Function Quick List

### change
- `ISliceToString([]int64{1, 2, 3}) // -> []string{"1","2","3"}`
- `SSlicetoISlice([]string{"1","2","3"}) // -> []int64{1,2,3}`

### config
- `ReadConfig() // -> map[string]string of key=value pairs from the config file`
- `GetConfig() // -> []byte of raw config file contents`
- `OverwriteConfig(map[string]string{"KEY":"VALUE"}) // merges into existing config and writes`
- `WriteConfig(map[string]string{"KEY":"VALUE"}) // -> true on success`

### dataverse
- `Authenticate() // retrieves and caches a new access token`
- `GetAccessToken() // requests a token using ClientID, ClientSecret, TenantID, Endpoint`
- `Request("accounts", "$select=name&$top=10") // GET request, returns []byte body`
- `Create("accounts", []byte(`{"name":"Acme"}`)) // POST request, returns []byte body`

### guid
- `New("00000000-0000-0000-0000-000000000000") // -> Guid{String:..., Valid:true}`
- `Empty() // -> Guid{String:"", Valid:false}`
- `Matches(guid.New("..."), guid.New("...")) // -> true if equal (case-insensitive)`
- `MatchesString(guid.New("..."), "...") // -> true if equal (case-insensitive)`

### llm
- `UploadFile("/path/to/file.pdf", "user_data") // uploads, extracts, structures -> []byte JSON`
- `UseUploadedFile("file_abc123") // fetches and structures contents for an existing file ID`
- `GetContents("file_abc123") // extracts a TOC from the uploaded file -> []byte JSON`
- `StructureContents(rawJSON) // cleans/validates TOC JSON -> []byte JSON`
- `SendRequest([]PromptMessage{{Role:"user", Content: []struct{Type, Text, Image, File string}{{Type:"input_text", Text:"Hello"}}}}, "gpt-5-2025-08-07") // low-level call to responses API`

### logging
- `InitLogs() // initializes log files and rotation (set AppName/ExecutableName first)`
- `InitVars() // sets ConfigPath/BaseLogsFolder (usually called by InitLogs)`
- `InitErrorLog("/tmp/errors-2025-01-01.log") // opens error log file`
- `InitChangeLog("/tmp/changes-2025-01-01.log") // opens change log file`
- `InitTraceLog("/tmp/trace-2025-01-01.log") // opens trace log file (no-op if !TraceDebug)`
- `SetLogsFolder("session-123") // switches logs to subfolder session-123/`
- `ErrorLog("something failed") // writes to error log`
- `ChangeLog("updated record", "ID-42") // writes to change log with ID`
- `TraceLog("debug details") // writes to trace log (only if TraceDebug)`
- `RetrieveLatestCaller("msg") // -> "time || (pkg) Func:Line || msg"`
- `PrintLogs("msg", 0) // prints as error (0=error,1=change,2=trace)`
- `RotateLogs(BaseLogsFolder) // removes logs older than retention`
- `LogOutdated("2025-1-02", 7) // -> true if older than 7 days`
- `ClearOutdatedLogs("/tmp/trace-2024-12-01.log") // deletes old log`
- `PanicErr(err) // panics if err != nil`
- `PrintErr(err) // logs error and continues`
- `ErrExists(err) // -> true if err != nil`
- `RetrieveErr(err) // -> err (and logs if not nil)`
- `defer WrapErr(file.Close) // logs error from defer and continues`
- `defer WrapPanic(file.Close) // logs error from defer and exits`
- `Panic("critical failure") // logs with stack and exits`
- `Assert(10, 10) // panics if types/values mismatch`

### net
- `SSHTunnel(client, "127.0.0.1:8080", "remote.host:80"); defer ln.Close() // creates TCP tunnel via SSH`
- `LoadDefaultPrivateKeys() // reads ~/.ssh/id_ed25519 or id_rsa`
- `CreateShellStream(&w, "ls", "-la") // streams shell output to HTTP client`
- `CreateInternalStream(&w) // returns writers streaming to client`

### os
- `HideConsole(cmd) // Windows-only: hides window for process`
- `Run("echo hello") // runs shell command, returns combined output and ok`
- `EnsurePath("/tmp/app/config.ini") // creates parent dirs if needed`
- `SetWindowFullscreen("Untitled - Notepad") // Windows-only: makes window fullscreen`
- `GetScreenSize() // Windows-only: screen width/height`
- `GetSystemMetrics(SM_CXSCREEN) // Windows-only: metric value`
- `FindWindowByTitle("Untitled - Notepad") // Windows-only: window handle`
- `SetWindowPos(hwnd, 0, 0, 0, w, h, SWP_SHOWWINDOW|SWP_FRAMECHANGED) // Windows-only: moves/resizes window`
- `GetAllActiveWindows() // Windows-only: []Window of visible top-level windows`
- `type Window struct { Title, FullTitle string; Handle uintptr; Process uint32; Executable string }` // Windows-only: Window information struct
- `const SM_CXSCREEN, SM_CYSCREEN, SWP_SHOWWINDOW, SWP_FRAMECHANGED` // Windows-only

### test
- `BenchmarkFunctions(fn1, fn2, 1000) // logs total and per-call ms for two funcs`

### time
- `UnixSince(ts, 60) // -> true if now - ts > 60 seconds`
- `UnixUntil(ts, 60) // -> true if now - ts < 60 seconds`
- `ConvertUnixToTimestamp(1710000000) // -> "YYYY-MM-DD HH:MM:SS.mmm"`
- `GetTime() // -> "3:04PM" (local)`
- `GetUnixTime() // -> int64 seconds UTC`
- `GetTimestamp() // -> "YYYY-MM-DD HH:MM:SS.mmm" (UTC)`
- `GetFullDateTime() // -> "D/M/YY h:mmam/pm" (local)`
- `GetDay() // -> "YYYY-M-D" (UTC)`
- `GetDateArray() // -> []int{day, month, year} (local)`
