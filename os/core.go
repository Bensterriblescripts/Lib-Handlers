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
	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell.exe", "-Command", command)
		if out, err := ErrorExists(cmd.CombinedOutput()); err {
			return string(out), false
		} else {
			return string(out), true
		}
	} else {
		cmd := exec.Command("bash", "-c", command)
		if out, err := ErrorExists(cmd.CombinedOutput()); err {
			return string(out), false
		} else {
			return string(out), true
		}
	}
}
func EnsurePath(path string) bool {
	if ErrExists(os.MkdirAll(filepath.Dir(path), 0755)) {
		ErrorLog("Failed to create directory: " + path)
		return false
	}
	return true
}
