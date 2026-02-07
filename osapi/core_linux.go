package osapi

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

// Run executes a PowerShell command and returns its output and success status.
//
// Example:
//
//	out, ok := osapi.Run("Get-Date")
func Bash(command string) (string, bool) {
	cmd := exec.Command("bash", "-c", command)
	var stdOutBuf bytes.Buffer
	var stdErrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(&stdOutBuf, os.Stdout, TraceLogFile)
	cmd.Stderr = io.MultiWriter(&stdErrBuf, os.Stderr, TraceLogFile)

	if ErrExists(cmd.Run()) {
		ErrorLog("Failed to run command: " + command)
		return stdErrBuf.String(), false
	} else {
		return stdOutBuf.String(), true
	}
}

// EnsurePath creates all directories needed for the given path.
//
// Example:
//
//	ok := osapi.EnsurePath("C:\\Local\\Config\\app.ini")
func EnsurePath(path string) bool {
	if ErrExists(os.MkdirAll(filepath.Dir(path), 0755)) {
		ErrorLog("Failed to create directory: " + path)
		return false
	}
	return true
}

// GetFileSize returns the file size in bytes.
//
// Example:
//
//	size := osapi.GetFileSize("C:\\Local\\Config\\app.ini")
func GetFileSize(path string) int64 {
	if info, err := os.Stat(path); err != nil {
		ErrorLog("Failed to get file size")
		return 0
	} else {
		if info.IsDir() {
			ErrorLog("Path is a directory, not a file " + path)
			return 0
		}

		return info.Size()
	}
}

// AddToLocalSoftware copies the current executable to C:\\Local\\Software.
//
// Example:
//
//	ok := osapi.AddToLocalSoftware()
func AddToLocalSoftware() bool {
	if currentExe, err := os.Executable(); err == nil {
		if currentExe == "C:\\Local\\Software\\"+AppName+".exe" {
			TraceLog("The executable is already running from this path...")
			return true
		}

		if _, err := os.Stat("C:\\Local\\Software\\" + AppName + ".exe"); err == nil { // Check if it's already there
			if GetFileSize(currentExe) == GetFileSize("C:\\Local\\Software\\"+AppName+".exe") { // They must be different sizes to bother copying
				TraceLog("No changes were made to the executable, leaving it as is...")
				return true
			}
		}

		if !EnsurePath("C:\\Local\\Software\\") {
			ErrorLog("Failed to ensure the path exists: C:/Local/Software/")
			return false
		}

		out, success := Bash(`Copy-Item -Path ` + currentExe + ` -Destination "C:\Local\Software\` + AppName + `.exe" -Force`)
		if !success {
			ErrorLog("Failed to copy self to software: " + out)
			return false
		}
		TraceLog("Copied the current executable into C:/Local/Software/" + AppName + ".exe")
		return true

	} else {
		ErrorLog("Failed to get executable path: " + err.Error())
		return false
	}
}
