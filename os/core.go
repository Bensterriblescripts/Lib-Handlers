package os

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

func Run(command string) string {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell.exe", "-Command", command)
		if out, err := ErrorExists(cmd.CombinedOutput()); err {
			return ""
		} else {
			return string(out)
		}
	} else {
		cmd := exec.Command("bash", "-c", command)
		if out, err := ErrorExists(cmd.CombinedOutput()); err {
			return ""
		} else {
			return string(out)
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
