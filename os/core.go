package osapi

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

func HideConsole(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
func Run(command string) (string, bool) {
	cmd := exec.Command("bash", "-c", command)

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
