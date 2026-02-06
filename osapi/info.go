package osapi

import "unsafe"

// GetScreenSize returns the primary screen width and height.
//
// Example:
//
//	w, h := osapi.GetScreenSize()
func GetScreenSize() (width, height int32) {
	w := GetSystemMetrics(SM_CXSCREEN)
	h := GetSystemMetrics(SM_CYSCREEN)
	return w, h
}
// GetSystemMetrics returns a system metric for the provided index.
//
// Example:
//
//	width := osapi.GetSystemMetrics(osapi.SM_CXSCREEN)
func GetSystemMetrics(index int32) int32 {
	r, _, _ := procGetSystemMetrics.Call(uintptr(index))
	return int32(r)
}
// GetMessage wraps the Win32 GetMessage call.
//
// Example:
//
//	result := osapi.GetMessage(&msg, 0, 0, 0)
func GetMessage(msg *MSG, hwnd uintptr, msgFilterMin uint32, msgFilterMax uint32) int32 {
	r, _, _ := procGetMessageW.Call(
		uintptr(unsafe.Pointer(msg)),
		hwnd,
		uintptr(msgFilterMin),
		uintptr(msgFilterMax),
	)
	return int32(r)
}
