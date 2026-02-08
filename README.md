# Lib-Handlers

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.25.6-blue?logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)]()
[![Status](https://img.shields.io/badge/Status-Experimental-orange)]()

> A personal utility library for Go, heavily tailored to my own workflow. Documented here in case it's referenced elsewhere.
> Platform-specific functions (Windows API, registry, keylogger, etc.) use build tags and are excluded on Linux.

---

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Packages](#packages)
  - [logging](#logging) - Error handling, log files, rotation
  - [config](#config) - INI-style key=value config files
  - [dataverse](#dataverse) - Microsoft Dataverse OData API client
  - [guid](#guid) - GUID type with validation and JSON support
  - [llm](#llm) - OpenAI Responses API wrapper
  - [mutate](#mutate) - Slice and string transformations
  - [network](#network) - SSH tunnelling and HTTP streaming
  - [osapi](#osapi) - Shell commands, file helpers, Windows API
  - [test](#test) - Function benchmarking
  - [time](#time) - Timestamps, comparisons, formatting

---

## Installation

```sh
go get github.com/Bensterriblescripts/Lib-Handlers@latest
```

---

## Quick Start

Set `logging.AppName` and call `InitLogs()` before using the library. Log files are created automatically under `C:\Local\Logs\<AppName>\` (Windows) or `~/local/Logs/<AppName>/` (Linux).

```go
package main

import "github.com/Bensterriblescripts/Lib-Handlers/logging"

func main() {
    logging.AppName = "MyApp"
    logging.InitLogs()
    defer logging.ErrorLogFile.Close()
    defer logging.ChangeLogFile.Close()
    defer logging.TraceLogFile.Close()
}
```

Optional variables to set before `InitLogs()`:

| Variable | Type | Default | Description |
|---|---|---|---|
| `logging.AppName` | `string` | `"Default"` | Application name, used for log folder and config paths |
| `logging.LoggingPath` | `string` | `""` | Override the default log directory |
| `logging.TraceDebug` | `bool` | `true` | Enable trace-level logging |
| `logging.ConsoleLogging` | `bool` | `true` | Mirror log output to stdout |
| `logging.TraceLogRotation` | `int64` | `3` | Days before trace logs are rotated |
| `logging.PriorityLogRotation` | `int64` | `14` | Days before error/change logs are rotated |

---

## Packages

### logging

Error handling wrappers and structured log file management with automatic rotation.

**Log Writers**

```go
logging.ErrorLog("connection refused")              // writes to error log
logging.ChangeLog("updated profile", "user-42")     // writes to change log with ID
logging.TraceLog("starting sync job")               // writes to trace log (only if TraceDebug)
```

**Error Handling**

Functions suffixed with `Err` accept `(error)`. Functions suffixed with `Error` accept `(T, error)` using generics.

```go
// (error) variants
logging.PanicErr(err)                     // log + os.Exit if err != nil
logging.PrintErr(err)                     // log + continue
logging.ErrExists(err)                    // -> true if err != nil (logs it)

// (T, error) generic variants
val := logging.PanicError(riskyFunc())    // log + os.Exit if err != nil, else returns T
val := logging.PrintError(riskyFunc())    // log + continue, returns T
val, failed := logging.ErrorExists(riskyFunc()) // -> (T, true) if err != nil

// Defer helpers
defer logging.WrapErr(file.Close)         // log error from deferred call
defer logging.WrapPanic(file.Close)       // log + os.Exit from deferred call

// Utilities
logging.Panic("fatal: out of memory")     // log with stack trace + os.Exit
logging.Assert(expected, actual)           // panic if values/types differ
```

**Log Management**

```go
logging.InitLogs()                               // initialise all log files and start rotation
logging.SetLogsFolder("session-123")             // switch logs to a subfolder
logging.RotateLogs(logging.BaseLogsFolder)       // manually trigger log rotation
logging.RemoveLog("/path/to/old.log")            // delete a specific log file
logging.RetrieveLatestCaller("msg")              // -> "time || (pkg) Func:Line || msg"
logging.PrintLogs("msg", 0)                      // write to a log stream (0=error, 1=change, 2=trace)
```

---

### config

INI-style `key=value` config file stored at `C:\Local\Config\<AppName>.ini` by default.

```go
settings := config.ReadConfig()                              // -> map[string]string
config.Current["theme"] = "dark"
config.Write()                                               // persist config.Current to disk

config.WriteSetting("username", "alice")                     // write a single key
config.WriteSettings(map[string]string{"env": "prod"})       // merge multiple keys into config
```

| Variable | Description |
|---|---|
| `config.ConfigPath` | Override the default config file path |
| `config.Current` | The current in-memory config (populated by `ReadConfig`) |

---

### dataverse

Microsoft Dataverse / Dynamics 365 OData v9.2 API client with OAuth2 client-credentials auth.

**Setup**

```go
dataverse.ClientID     = guid.New("...")
dataverse.ClientSecret = "..."
dataverse.TenantID     = guid.New("...")
dataverse.Endpoint     = "https://org.crm6.dynamics.com"

dataverse.Authenticate()  // retrieves and caches an OAuth2 access token
```

**Querying**

```go
body := dataverse.Query("GET", "contacts", "$select=fullname&$top=10")
body := dataverse.Retrieve("contacts", "lastname eq 'Smith'", "fullname,emailaddress1", "lastname asc")
body := dataverse.RetrieveByID("contacts", "00000000-0000-0000-0000-000000000000", "fullname")
body := dataverse.RetrieveNext(nextLink)
```

**Writing**

```go
ok := dataverse.Create("contacts", map[string]interface{}{"fullname": "Ada Lovelace"})
ok := dataverse.Update("contacts", "record-guid", map[string]interface{}{"jobtitle": "Engineer"})
```

**Auth Helpers**

```go
dataverse.IsAuthenticated()       // -> true if the token is present and not expired
dataverse.EnsureAuthenticated()   // re-authenticates if needed, returns success
```

| Variable | Default | Description |
|---|---|---|
| `dataverse.VerboseLogging` | `true` | Log request/response details to trace log |
| `dataverse.AllowAnnotations` | `true` | Include OData annotations in responses |
| `dataverse.MaxPageSize` | `5000` | Max records per page (0 = Dataverse default) |

---

### guid

GUID type with validation, normalisation (lowercase, trimmed), and JSON unmarshalling.

```go
id := guid.New("7d444840-9dc0-11d1-b245-5ffdce74fad2")  // -> Guid{String: "...", Valid: true}
empty := guid.New()                                       // -> Guid{String: "0000...", Valid: false}

ok := guid.MatchesString(id, "7D444840-...")              // case-insensitive comparison
// Two Guid values can also be compared directly: a.String == b.String
```

Implements `json.Unmarshaler` so it works automatically with `json.Unmarshal`.

---

### llm

OpenAI Responses API wrapper for file uploads and chat-style prompts.

```go
llm.OpenAIApiKey = "sk-..."

resp := llm.UploadFile("/path/to/file.pdf", "user_data")     // upload + extract + structure -> JSON
resp := llm.UseUploadedFile("file_abc123")                    // structure an already-uploaded file
resp := llm.GetContents("file_abc123")                        // extract table of contents -> JSON
resp := llm.StructureContents(rawJSON)                        // clean/validate TOC JSON

resp := llm.SendRequest([]llm.PromptMessage{
    {Role: "user", Content: []struct{ Type, Text, Image, File string }{
        {Type: "input_text", Text: "Hello"},
    }},
}, "gpt-5-2025-08-07")
```

---

### mutate

Slice and string transformation helpers.

```go
mutate.ISliceToString([]int64{1, 2, 3})          // -> []string{"1", "2", "3"}
mutate.SSlicetoISlice([]string{"1", "2", "3"})   // -> []int64{1, 2, 3}
mutate.Capitalise("hello world")                  // -> "Hello world"
```

---

### network

SSH tunnelling and HTTP streaming utilities.

```go
ln := network.SSHTunnel(client, "127.0.0.1:8080", "remote.host:80")
defer ln.Close()

signer := network.LoadDefaultPrivateKeys()   // reads ~/.ssh/id_ed25519 or id_rsa

output := network.CreateShellStream(&w, "ls", "-la")   // stream shell output to HTTP client
stdout, stderr := network.CreateInternalStream(&w)      // get writers that stream to HTTP client
```

---

### osapi

Shell execution, file helpers, and Windows-specific API calls (window management, keylogger, hotkeys, registry, system info).

**Shell & File Helpers** (cross-platform)

```go
out, ok := osapi.PowerShell("Get-Date")       // Windows
out, ok := osapi.Bash("date")                 // Linux

ok := osapi.EnsurePath("C:\\Local\\Config\\app.ini")   // create parent directories
size := osapi.GetFileSize("C:\\file.txt")               // file size in bytes
ok := osapi.AddToLocalSoftware()                        // copy running exe to C:\Local\Software\
```

**Window Management** (Windows only)

```go
hwnd := osapi.GetWindowByTitle("Notepad")
title := osapi.GetWindowTitle(hwnd)
state := osapi.GetWindowState(hwnd)                         // "normal", "minimized", "maximized"
rect := osapi.GetWindowRect(hwnd)                           // -> RECT{Left, Top, Right, Bottom}
windows := osapi.GetAllActiveWindows()                      // -> []Window

osapi.SetBorderlessWindow(hwnd)                             // remove title bar and borders
osapi.SetWindowWindowed(hwnd)                               // restore windowed mode
osapi.SetWindowMinimised(hwnd)                              // minimise
osapi.SetFocus(hwnd)                                        // bring to front
osapi.SetVisible(hwnd)                                      // show + restore
osapi.SetWindowPos(hwnd, 0, 0, 0, 1920, 1080, osapi.SWP_SHOWWINDOW|osapi.SWP_FRAMECHANGED)

width, height := osapi.GetScreenSize()
metric := osapi.GetSystemMetrics(osapi.SM_CXSCREEN)
monitorRect := osapi.GetMonitorByWindow(hwnd)               // -> RECT of the monitor containing hwnd
```

**Keylogger & Hotkeys** (Windows only)

```go
osapi.AddHotkey("ctrl", "f1", func() { fmt.Println("hotkey pressed") })
osapi.StartKeylogger()   // blocks; processes hotkeys and optionally logs keystrokes
osapi.StopKeylogger()
```

| Variable | Description |
|---|---|
| `osapi.LogKeys` | When `true`, keystrokes are logged (default `false`) |
| `osapi.Hotkeys` | Registered hotkeys slice |

**Registry** (Windows only)

```go
osapi.RunExeAtLogon("MyApp", "C:\\Local\\Software\\MyApp.exe")   // add to HKCU Run key
```

**System Info** (Windows only)

```go
osapi.CurrentProcessUsage()   // logs current process CPU and memory usage
```

**Types & Constants** (Windows only)

```go
type Window struct {
    Title, FullTitle string
    Handle           uintptr
    Process          uint32
    Executable       string
    WindowState      string
    OriginalRect     RECT
    MonitorInfo      RECT
}

type RECT struct{ Left, Top, Right, Bottom int32 }
type MONITORINFO struct{ CbSize uint32; RcMonitor, RcWork RECT; DwFlags uint32 }
type MSG struct{ Hwnd uintptr; Message, WParam, LParam uintptr; Time uint32; Pt POINT; LPrivate uint32 }
type POINT struct{ X, Y int32 }
type Hotkey struct{ ID int; Mod, Key uintptr; Callback func(); Active bool }
```

Constants: `SM_CXSCREEN`, `SM_CYSCREEN`, `SWP_FRAMECHANGED`, `SWP_SHOWWINDOW`, `GWL_STYLE`, `WS_POPUP`, `WS_CAPTION`, `WS_THICKFRAME`, `WS_MINIMIZEBOX`, `WS_MAXIMIZEBOX`, `WS_SYSMENU`, `WS_OVERLAPPEDWINDOW`, `MONITOR_DEFAULTTONEAREST`, `SW_SHOW`, `SW_RESTORE`, `SW_SHOWMAXIMIZED`, `SW_MINIMIZE`

---

### test

Simple function benchmarking.

```go
test.BenchmarkFunctions(fn1, fn2, 1000)   // compare two functions over N iterations
test.BenchmarkFunction(fn, 1000)           // benchmark a single function over N iterations
```

---

### time

Timestamp generation, comparison, and formatting.

**Comparison**

```go
time.UnixSince(ts, 60)    // -> true if more than 60 seconds have passed since ts
time.UnixUntil(ts, 60)    // -> true if fewer than 60 seconds have passed since ts
```

**Retrieval**

```go
time.GetTime()             // -> "3:04PM" (local)
time.GetUnixTime()         // -> int64 seconds (UTC)
time.GetTimestamp()         // -> "YYYY-MM-DD HH:MM:SS.mmm" (UTC)
time.GetFullDateTime()     // -> "D/M/YY h:mmam/pm" (local)
time.GetDay()              // -> "YYYY-M-D" (UTC)
time.GetDateArray()        // -> []int{day, month, year} (local)
```

**Conversion**

```go
time.ConvertUnixToTimestamp(1710000000)   // -> "YYYY-MM-DD HH:MM:SS.mmm"
```
