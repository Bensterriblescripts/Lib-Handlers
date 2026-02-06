package osapi

import (
	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

// RunExeAtLogon registers an executable to run at user logon.
//
// Example:
//
//	ok := osapi.RunExeAtLogon("MyApp", "C:\\Local\\Software\\MyApp.exe")
func RunExeAtLogon(name string, path string) bool {
	out, success := Run(`
		New-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Run" -Name "` + name + `" -Value '"` + path + `"' -PropertyType String -Force
	`)
	if !success {
		ErrorLog("Failed to create task: " + name)
		ErrorLog(out)
		return false
	} else {
		TraceLog("Created new executable task at logon: " + name)
		return true
	}
}
