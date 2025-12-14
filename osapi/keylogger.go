package osapi

import (
	"runtime"
	"strconv"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

var LogKeys = true

const (
	MOD_CONTROL = 0x0002
	MOD_SHIFT   = 0x0004

	hotkeyID = 1
	L        = 0x4C
)

func StartHotkeyService(modifiers uint32, key uint32, Callback func()) {
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		if !RegisterHotKey(0, hotkeyID, modifiers, key) {
			ErrorLog("RegisterHotKey(" + strconv.FormatUint(uint64(modifiers), 16) + "+" + strconv.FormatUint(uint64(key), 16) + ") failed")
			return
		}
		defer UnregisterHotKey(0, hotkeyID)

		for LogKeys {
			var msg MSG
			r := GetMessage(&msg, 0, 0, 0)
			if r <= 0 {
				break
			}
			if msg.Message == WM_HOTKEY && msg.WParam == uintptr(hotkeyID) {
				TraceLog("Triggered Hotkey")
				Callback()
			}
		}
	}()
}

func RegisterHotKey(hwnd uintptr, id int32, modifiers uint32, vk uint32) bool {
	r, _, err := procRegisterHotKey.Call(
		hwnd,
		uintptr(id),
		uintptr(modifiers),
		uintptr(vk),
	)
	if r == 0 {
		ErrorLog("RegisterHotKey failed: " + err.Error())
		return false
	}
	return true
}

func UnregisterHotKey(hwnd uintptr, id int32) bool {
	r, _, err := procUnregisterHotKey.Call(
		hwnd,
		uintptr(id),
	)
	if r == 0 {
		ErrorLog("RegisterHotKey failed: " + err.Error())
		return false
	}
	return true
}
