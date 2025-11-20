package osapi

import (
	"fmt"
	"syscall"
	"unsafe"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
	"golang.org/x/sys/windows"
)

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	version  = windows.NewLazySystemDLL("version.dll")

	/* Enumeration */
	procEnumWindows = user32.NewProc("EnumWindows")

	/* Get Window Information */
	procGetWindowText       = user32.NewProc("GetWindowTextW")
	procGetWindowTextLength = user32.NewProc("GetWindowTextLengthW")
	procIsWindowVisible     = user32.NewProc("IsWindowVisible")
	procFindWindowW         = user32.NewProc("FindWindowW")
	procGetSystemMetrics    = user32.NewProc("GetSystemMetrics")

	/* Process */
	procGetWindowThreadProcessId   = user32.NewProc("GetWindowThreadProcessId")
	procQueryFullProcessImageNameW = kernel32.NewProc("QueryFullProcessImageNameW")

	/* Executable Information */
	procGetFileVersionInfoSizeW = version.NewProc("GetFileVersionInfoSizeW")
	procGetFileVersionInfoW     = version.NewProc("GetFileVersionInfoW")
	procVerQueryValueW          = version.NewProc("VerQueryValueW")

	/* Change States */
	procSetWindowPos = user32.NewProc("SetWindowPos")
)

const (
	SM_CXSCREEN      = 0 // width of primary monitor
	SM_CYSCREEN      = 1 // height of primary monitor
	SWP_FRAMECHANGED = 0x0020
	SWP_SHOWWINDOW   = 0x0040
)

/* Window Management */
func SetWindowFullscreen(windowname string) {
	hwnd := FindWindowByTitle(windowname)
	if hwnd == 0 {
		return
	}
	screenWidth, screenHeight := GetScreenSize()

	SetWindowPos(hwnd, 0, 0, 0, screenWidth, screenHeight, SWP_SHOWWINDOW|SWP_FRAMECHANGED)
}
func SetWindowPos(hwnd uintptr, hwndInsertAfter uintptr, x, y, cx, cy int32, flags uint32) bool {
	r, _, _ := procSetWindowPos.Call(
		hwnd,
		hwndInsertAfter,
		uintptr(x),
		uintptr(y),
		uintptr(cx),
		uintptr(cy),
		uintptr(flags),
	)
	return r != 0
}

/* Window Information */
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
func GetScreenSize() (width, height int32) {
	w := GetSystemMetrics(SM_CXSCREEN)
	h := GetSystemMetrics(SM_CYSCREEN)
	return w, h
}
func GetSystemMetrics(index int32) int32 {
	r, _, _ := procGetSystemMetrics.Call(uintptr(index))
	return int32(r)
}

type Window struct {
	Title      string
	FullTitle  string
	Handle     uintptr
	Process    uint32
	Executable string
}

var activeWindows []Window

func enumWindowsCallback(hwnd uintptr, _ uintptr) uintptr {
	if visible, _, _ := procIsWindowVisible.Call(hwnd); visible == 0 {
		return 1
	}

	var window Window
	window.Handle = hwnd

	// Process ID
	if _, _, e := procGetWindowThreadProcessId.Call(
		hwnd,
		uintptr(unsafe.Pointer(&window.Process)),
	); window.Process == 0 {
		if e != nil && e.Error() != syscall.Errno(0).Error() {
			ErrorLog(fmt.Sprintf("GetWindowThreadProcessId failed for hwnd=0x%x: %v", hwnd, e))
		} else {
			TraceLog(fmt.Sprintf("Skipping window with PID 0: hwnd=0x%x", hwnd))
		}
		return 1
	}

	// Window Title/FullTitle
	window.Title = getWindowTitleByHandle(hwnd)
	if window.Title == "" {
		return 1
	}
	window.FullTitle = window.Title

	// Executable path
	exePath, err := getProcessImagePath(window.Process)
	if err != nil {
		ErrorLog(fmt.Sprintf("getProcessImagePath failed for PID %d: %v", window.Process, err))
	} else {
		window.Executable = exePath
	}

	// Windows Title - If Description is Available
	if window.Executable != "" {
		if desc, ferr := getFileDescriptionByPath(window.Executable); ferr != nil {
			ErrorLog(fmt.Sprintf("getFileDescriptionByPath(%q) failed: %v", window.Executable, ferr))
		} else if desc != "" {
			window.Title = desc
		}
	}

	TraceLog(fmt.Sprintf("Found Window: hwnd=0x%x, title=%q", hwnd, window.Title))
	activeWindows = append(activeWindows, window)
	return 1
}

func getProcessImagePath(pid uint32) (path string, err error) {
	const PROCESS_QUERY_LIMITED_INFORMATION = 0x1000

	h, err := windows.OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return "", err
	}
	defer func() {
		err := windows.CloseHandle(h)
		if err != nil {
			ErrorLog(fmt.Sprintf("CloseHandle failed for pid=%d: %v", pid, err))
		}
	}()

	buf := make([]uint16, 260)
	size := uint32(len(buf))

	r0, _, e := procQueryFullProcessImageNameW.Call(
		uintptr(h),
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
	)
	if r0 == 0 {
		if e != nil && e.Error() != syscall.Errno(0).Error() {
			err = e
			return "", err
		}
		err = fmt.Errorf("QueryFullProcessImageNameW returned 0 without extended error")
		return "", err
	}

	path = windows.UTF16ToString(buf[:size])
	return path, nil
}
func getFileDescriptionByPath(path string) (desc string, err error) {
	if path == "" {
		return "", nil
	}

	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return "", err
	}

	var handle uint32
	r0, _, callErr := procGetFileVersionInfoSizeW.Call(
		uintptr(unsafe.Pointer(p)),
		uintptr(unsafe.Pointer(&handle)),
	)
	if r0 == 0 {
		TraceLog(fmt.Sprintf("No version info available for %q", path))
		return "", nil
	}
	size := uint32(r0)

	buf := make([]byte, size)
	r0, _, callErr = procGetFileVersionInfoW.Call(
		uintptr(unsafe.Pointer(p)),
		0,
		uintptr(size),
		uintptr(unsafe.Pointer(&buf[0])),
	)
	if r0 == 0 {
		if callErr != nil && callErr.Error() != syscall.Errno(0).Error() {
			err = callErr
			return "", err
		}
		err = fmt.Errorf("GetFileVersionInfoW returned 0 for %q", path)
		return "", err
	}

	var transPtr uintptr
	var transLen uint32
	transStringPtr, _ := syscall.UTF16PtrFromString(`\VarFileInfo\Translation`)
	if transStringPtr == nil {
		return "", fmt.Errorf("UTF16PtrFromString failed for `\\VarFileInfo\\Translation`")
	}
	r0, _, callErr = procVerQueryValueW.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&transStringPtr)),
		uintptr(unsafe.Pointer(&transPtr)),
		uintptr(unsafe.Pointer(&transLen)),
	)
	if r0 == 0 || transLen < 4 {
		if callErr != nil && callErr.Error() != syscall.Errno(0).Error() {
			TraceLog(fmt.Sprintf("VerQueryValueW(Translation) fallback for %q: %v", path, callErr))
		}
		return queryFileDescription(buf, 0x0409, 0x04B0)
	}

	lang := uint16(transPtr)
	codepage := uint16(transPtr + unsafe.Sizeof(lang))

	return queryFileDescription(buf, lang, codepage)
}

func queryFileDescription(buf []byte, lang, codepage uint16) (desc string, err error) {
	subBlock := fmt.Sprintf(`\StringFileInfo\%04x%04x\FileDescription`, lang, codepage)

	var valuePtr uintptr
	var valueLen uint32
	subBlockPtr, _ := syscall.UTF16PtrFromString(subBlock)
	if subBlockPtr == nil {
		return "", fmt.Errorf("UTF16PtrFromString failed for %q", subBlock)
	}
	r0, _, callErr := procVerQueryValueW.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(subBlockPtr)),
		uintptr(unsafe.Pointer(&valuePtr)),
		uintptr(unsafe.Pointer(&valueLen)),
	)
	if r0 == 0 || valueLen == 0 {
		if callErr != nil && callErr.Error() != syscall.Errno(0).Error() {
			err = callErr
			return "", err
		}
		err = fmt.Errorf("VerQueryValueW returned no data for %q", subBlock)
		return "", err
	}

	desc = windows.UTF16PtrToString((*uint16)(unsafe.Pointer(&valuePtr)))
	TraceLog(fmt.Sprintf("File description (%s): %q", subBlock, desc))
	return desc, nil
}
func getWindowTitleByHandle(hwnd uintptr) string {
	ret, _, callErr := procGetWindowTextLength.Call(hwnd)
	length := uint32(ret)
	if length == 0 {
		if callErr != nil && callErr.Error() != syscall.Errno(0).Error() {
			TraceLog(fmt.Sprintf("GetWindowTextLength failed for hwnd=0x%x: %v", hwnd, callErr))
		}
		return ""
	}
	buf := make([]uint16, length+1)

	_, _, callErr = procGetWindowText.Call(
		hwnd,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(length+1),
	)
	if callErr != nil && callErr.Error() != syscall.Errno(0).Error() {
		TraceLog(fmt.Sprintf("GetWindowText failed for hwnd=0x%x: %v", hwnd, callErr))
	}

	return windows.UTF16ToString(buf)
}

func GetAllActiveWindows() []Window {
	TraceLog("GetAllActiveWindows started")
	activeWindows = make([]Window, 0)

	cb := syscall.NewCallback(enumWindowsCallback)
	ret, _, err := procEnumWindows.Call(
		cb,
		0,
	)

	if ret == 0 {
		if err != nil && err.Error() != syscall.Errno(0).Error() {
			ErrorLog("Error while iterating through windows: " + err.Error())
		} else {
			ErrorLog("EnumWindows returned 0 without extended error")
		}
		return nil
	}

	TraceLog(fmt.Sprintf("GetAllActiveWindows finished, found %d windows", len(activeWindows)))
	return activeWindows
}
