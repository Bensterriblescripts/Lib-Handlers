package osapi

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	procSetWindowPos     = user32.NewProc("SetWindowPos")
	procFindWindowW      = user32.NewProc("FindWindowW")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")

	kernel32              = windows.NewLazySystemDLL("kernel32.dll")
	procGetCurrentProcess = kernel32.NewProc("GetCurrentProcess")
	procSetPriorityClass  = kernel32.NewProc("SetPriorityClass")
	procGetCurrentThread  = kernel32.NewProc("GetCurrentThread")
	procSetThreadPriority = kernel32.NewProc("SetThreadPriority")
)

const (
	SM_CXSCREEN      = 0 // width of primary monitor
	SM_CYSCREEN      = 1 // height of primary monitor
	SWP_FRAMECHANGED = 0x0020
	SWP_SHOWWINDOW   = 0x0040

	PROCESS_MODE_BACKGROUND_BEGIN = 0x00100000
	PROCESS_MODE_BACKGROUND_END   = 0x00200000
	THREAD_MODE_BACKGROUND_BEGIN  = 0x00010000
	THREAD_MODE_BACKGROUND_END    = 0x00020000
)

func SetWindowFullscreen(windowname string) {
	hwnd := FindWindowByTitle(windowname)
	if hwnd == 0 {
		return
	}
	screenWidth, screenHeight := GetScreenSize()

	SetWindowPos(hwnd, 0, 0, 0, screenWidth, screenHeight, SWP_SHOWWINDOW|SWP_FRAMECHANGED)
}

func GetScreenSize() (width, height int32) {
	w := GetSystemMetrics(SM_CXSCREEN)
	h := GetSystemMetrics(SM_CYSCREEN)
	return w, h
}
func GetSystemMetrics(index int32) int32 {
	r, _, _ := procGetSystemMetrics.Call(uintptr(index))
	return int32(r)
}

func FindWindowByTitle(title string) uintptr {
	t, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return 0
	}
	r, _, _ := procFindWindowW.Call(
		0,
		uintptr(unsafe.Pointer(t)),
	)

	return r
}
func SetWindowPos(hwnd uintptr, hwndInsertAfter uintptr, x, y, cx, cy int32, flags uint32) bool {
	r, _, _ := procSetWindowPos.Call(
		uintptr(hwnd),
		uintptr(hwndInsertAfter),
		uintptr(x),
		uintptr(y),
		uintptr(cx),
		uintptr(cy),
		uintptr(flags),
	)
	return r != 0
}
