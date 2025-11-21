package osapi

import (
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

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
