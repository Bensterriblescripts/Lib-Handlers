package osapi

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

func HideConsole(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
func Run(command string) (string, bool) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", command)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // No powershell window
	} else {
		cmd = exec.Command("bash", "-c", command)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), false
	}
	return string(out), true
}
func EnsurePath(path string) bool {
	if ErrExists(os.MkdirAll(filepath.Dir(path), 0755)) {
		ErrorLog("Failed to create directory: " + path)
		return false
	}
	return true
}
