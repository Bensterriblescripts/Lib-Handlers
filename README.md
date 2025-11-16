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
You must add a defer to close the files and any ssh handlers before the program closes. I am too stupid/lazy to figure out how to do this outside of the main function at this moment.

---

## Function Quick List

### change
- `change.ISliceToString([]int64{1, 2, 3}) // -> []string{"1","2","3"}`
- `change.SSlicetoISlice([]string{"1","2","3"}) // -> []int64{1,2,3}`

### config
- `config.ReadConfig() // map[string]string of key=value pairs from the config file`
- `config.GetConfig() // []byte of raw config file contents`
- `config.OverwriteConfig(map[string]string{"KEY":"VALUE"}) // merges into existing and writes`
- `config.WriteConfig(map[string]string{"KEY":"VALUE"}) // true on success`

### dataverse
- `dataverse.Authenticate() // retrieves and caches a new access token`
- `dataverse.GetAccessToken() // requests a token using ClientID, ClientSecret, TenantID, Endpoint`
- `dataverse.Request("accounts", "$select=name&$top=10") // GET request, returns []byte body`
- `dataverse.Create("accounts", []byte(`{"name":"Acme"}`)) // POST request, returns []byte body`

### guid
- `guid.New("00000000-0000-0000-0000-000000000000") // -> Guid{String:..., Valid:true}`
- `guid.Empty() // -> Guid{String:"", Valid:false}`
- `guid.Matches(guid.New("..."), guid.New("...")) // true if equal (case-insensitive)`
- `guid.MatchesString(guid.New("..."), "...") // true if equal (case-insensitive)`

### llm
- `llm.UploadFile("/path/to/file.pdf", "user_data") // uploads, extracts, structures -> []byte JSON`
- `llm.UseUploadedFile("file_abc123") // fetches and structures contents for an existing file id`
- `llm.GetContents("file_abc123") // extracts a TOC from the uploaded file -> []byte JSON`
- `llm.StructureContents(rawJSON) // cleans/validates the TOC JSON -> []byte JSON`
- `llm.SendRequest([]llm.PromptMessage{{Role:"user", Content: []struct{Type, Text, Image, File string}{{Type:"input_text", Text:"Hello"}}}}, "gpt-5-2025-08-07") // low-level call to responses API`

### logging
- `logging.InitLogs() // initialize log files and rotation (set AppName/ExecutableName first)`
- `logging.InitVars() // sets ConfigPath/BaseLogsFolder (usually called by InitLogs)`
- `logging.InitErrorLog("/tmp/errors-2025-01-01.log") // open error log file`
- `logging.InitChangeLog("/tmp/changes-2025-01-01.log") // open change log file`
- `logging.InitTraceLog("/tmp/trace-2025-01-01.log") // open trace log file (no-op if !TraceDebug)`
- `logging.SetLogsFolder("session-123") // switch logs to subfolder session-123/`
- `logging.ErrorLog("something failed") // write to error log`
- `logging.ChangeLog("updated record", "ID-42") // write to change log with ID`
- `logging.TraceLog("debug details") // write to trace log (only if TraceDebug)`
- `logging.RetrieveLatestCaller("msg") // -> "time || (pkg) Func:Line || msg"`
- `logging.PrintLogs("msg", 0) // print as error (0=error,1=change,2=trace)`
- `logging.RotateLogs(logging.BaseLogsFolder) // remove logs older than retention`
- `logging.LogOutdated("2025-1-02", 7) // -> true if older than 7 days`
- `logging.ClearOutdatedLogs("/tmp/trace-2024-12-01.log") // delete old log`
- `logging.PanicErr(err) // panic if err != nil`
- `logging.PrintErr(err) // log error, continue`
- `logging.ErrExists(err) // -> true if err != nil`
- `logging.RetrieveErr(err) // -> err (and logs if not nil)`
- `defer logging.WrapErr(file.Close) // log error from defer and continue`
- `defer logging.WrapPanic(file.Close) // log error from defer and exit`
- `logging.Panic("critical failure") // log with stack and exit`
- `logging.Assert(10, 10) // panic if types/values mismatch`

### net (package network)
- `ln := network.SSHTunnel(client, "127.0.0.1:8080", "remote.host:80"); defer ln.Close() // TCP tunnel via SSH`
- `signer := network.LoadDefaultPrivateKeys() // read ~/.ssh/id_ed25519 or id_rsa`
- `network.CreateShellStream(&w, "ls", "-la") // stream shell output to HTTP client`
- `traceOut, errorOut := network.CreateInternalStream(&w) // get writers streaming to client`

### os (package osapi)
- `osapi.HideConsole(cmd) // Windows-only: hide window for process`
- `out, ok := osapi.Run("echo hello") // run shell command, get combined output`
- `ok := osapi.EnsurePath("/tmp/app/config.ini") // create parent dirs if needed`

### os (Windows HWND helpers; package osapi)
- `osapi.SetWindowFullscreen("Untitled - Notepad") // Windows-only: make window fullscreen`
- `w, h := osapi.GetScreenSize() // -> screen width/height`
- `osapi.GetSystemMetrics(osapi.SM_CXSCREEN) // -> metric value`
- `hwnd := osapi.FindWindowByTitle("Untitled - Notepad") // -> window handle`
- `ok := osapi.SetWindowPos(hwnd, 0, 0, 0, w, h, osapi.SWP_SHOWWINDOW|osapi.SWP_FRAMECHANGED) // move/resize`

### test
- `test.BenchmarkFunctions(fn1, fn2, 1000) // log total and per-call ms for two funcs`

### time
- `time.UnixSince(ts, 60) // -> true if now - ts > 60 seconds`
- `time.UnixUntil(ts, 60) // -> true if now - ts < 60 seconds`
- `time.ConvertUnixToTimestamp(1710000000) // -> "YYYY-MM-DD HH:MM:SS.mmm"`
- `time.GetTime() // -> "3:04PM" (local)`
- `time.GetUnixTime() // -> int64 seconds UTC`
- `time.GetTimestamp() // -> "YYYY-MM-DD HH:MM:SS.mmm" (UTC)`
- `time.GetFullDateTime() // -> "D/M/YY h:mmam/pm" (local)`
- `time.GetDay() // -> "YYYY-M-D" (UTC)`
- `time.GetDateArray() // -> []int{day, month, year} (local)`