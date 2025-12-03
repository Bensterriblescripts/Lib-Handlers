package osapi

import (
	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

func RunExecutableAtLogonRegistry(name string, path string) {
	out, success := Run(`
		New-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Run" -Name "` + name + `" -Value "` + path + `" -PropertyType String -Force
	`)
	if !success {
		ErrorLog("Failed to create task: " + name)
		ErrorLog(out)
	} else {
		TraceLog("Created new executable task at logon: " + name)
	}
}
