package osapi

import (
	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
	"golang.org/x/sys/windows/registry"
)

// RunExeAtLogon registers an executable to run at user logon.
//
// Example:
//
//	ok := osapi.RunExeAtLogon("MyApp", "C:\\Local\\Software\\MyApp.exe")
func RunExeAtLogon(name string, path string) bool {
	key, existed, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		ErrorLog("Failed to create registry key: " + err.Error())
		return false
	}
	if existed {
		TraceLog("Registry key already exists, overwriting value")
	}
	defer key.Close()

	key.SetStringValue(name, path)
	TraceLog("Created new executable task at logon: " + name)
	return true
}
