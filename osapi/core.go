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
