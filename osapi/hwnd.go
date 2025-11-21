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
	procSetWindowPos  = user32.NewProc("SetWindowPos")
	procSetWindowLong = user32.NewProc("SetWindowLongW")
	procGetWindowLong = user32.NewProc("GetWindowLongW")
)

const (
	SM_CXSCREEN             = 0 // width of primary monitor
	SM_CYSCREEN             = 1 // height of primary monitor
	SWP_FRAMECHANGED        = 0x0020
	SWP_SHOWWINDOW          = 0x0040
	SWP_NOSIZE              = 0x0001
	SWP_NOMOVE              = 0x0002
	SWP_NOZORDER            = 0x0004
	SWP_NOACTIVATE          = 0x0010
	SWP_NOOWNERZORDER       = 0x0200
	GWL_STYLE         int32 = -16
	GWL_EXSTYLE       int32 = -20
	WS_POPUP                = 0x80000000
	WS_VISIBLE              = 0x10000000
	WS_CLIPSIBLINGS         = 0x20000000
	WS_CLIPCHILDREN         = 0x40000000

	// Standard style bits to remove for borderless
	WS_CAPTION     = 0x00C00000
	WS_THICKFRAME  = 0x00040000
	WS_MINIMIZEBOX = 0x00020000
	WS_MAXIMIZEBOX = 0x00010000
	WS_SYSMENU     = 0x00080000

	// Extended style bits to remove for borderless
	WS_EX_DLGMODALFRAME = 0x00000001
	WS_EX_CLIENTEDGE    = 0x00000200
	WS_EX_STATICEDGE    = 0x00020000
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
	Style      uint32
}

var activeWindows []Window

func enumWindowsCallback(hwnd uintptr, _ uintptr) uintptr {
	// Visible windows only
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
		if e != nil && e != syscall.Errno(0) {
			ErrorLog(fmt.Sprintf("GetWindowThreadProcessId failed for hwnd=0x%x: %v", hwnd, e))
		} else {
			TraceLog(fmt.Sprintf("Skipping window with PID 0: hwnd=0x%x", hwnd))
		}
		return 1
	}

	// Window title
	window.Title = GetWindowTitleByHandle(hwnd)
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

	// Description from file version info if we have a path
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

// getProcessImagePath returns the full path of the process executable.
func getProcessImagePath(pid uint32) (path string, err error) {

	// Needs PROCESS_QUERY_LIMITED_INFORMATION – works for 64-bit targets enumerating 64-bit processes
	const PROCESS_QUERY_LIMITED_INFORMATION = 0x1000

	h, err := windows.OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return "", err
	}

	// Make sure we always close the handle, and log if CloseHandle fails.
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
		if e != nil && e != syscall.Errno(0) {
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
		// No version info – not necessarily an error
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
		if callErr != nil && callErr != syscall.Errno(0) {
			err = callErr
			return "", err
		}
		err = fmt.Errorf("GetFileVersionInfoW returned 0 for %q", path)
		return "", err
	}

	var transPtr uintptr
	var transLen uint32
	r0, _, callErr = procVerQueryValueW.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(`\VarFileInfo\Translation`))),
		uintptr(unsafe.Pointer(&transPtr)),
		uintptr(unsafe.Pointer(&transLen)),
	)
	if r0 == 0 || transLen < 4 {
		if callErr != nil && callErr != syscall.Errno(0) {
			TraceLog(fmt.Sprintf("VerQueryValueW(Translation) fallback for %q: %v", path, callErr))
		}
		// Fallback to US English / Unicode
		return queryFileDescription(buf, 0x0409, 0x04B0)
	}

	lang := *(*uint16)(unsafe.Pointer(transPtr))
	codepage := *(*uint16)(unsafe.Pointer(transPtr + unsafe.Sizeof(lang)))

	return queryFileDescription(buf, lang, codepage)
}

func queryFileDescription(buf []byte, lang, codepage uint16) (desc string, err error) {
	subBlock := fmt.Sprintf(`\StringFileInfo\%04x%04x\FileDescription`, lang, codepage)

	var valuePtr uintptr
	var valueLen uint32
	r0, _, callErr := procVerQueryValueW.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subBlock))),
		uintptr(unsafe.Pointer(&valuePtr)),
		uintptr(unsafe.Pointer(&valueLen)),
	)
	if r0 == 0 || valueLen == 0 {
		if callErr != nil && callErr != syscall.Errno(0) {
			err = callErr
			return "", err
		}
		err = fmt.Errorf("VerQueryValueW returned no data for %q", subBlock)
		return "", err
	}

	desc = windows.UTF16PtrToString((*uint16)(unsafe.Pointer(valuePtr)))
	TraceLog(fmt.Sprintf("File description (%s): %q", subBlock, desc))
	return desc, nil
}

func GetWindowTitleByHandle(hwnd uintptr) string {
	ret, _, callErr := procGetWindowTextLength.Call(hwnd)
	length := uint32(ret)
	if length == 0 {
		if callErr != nil && callErr != syscall.Errno(0) {
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
	if callErr != nil && callErr != syscall.Errno(0) {
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
		if err != nil && err != syscall.Errno(0) {
			ErrorLog("Error while iterating through windows: " + err.Error())
		} else {
			ErrorLog("EnumWindows returned 0 without extended error")
		}
		return nil
	}

	TraceLog(fmt.Sprintf("GetAllActiveWindows finished, found %d windows", len(activeWindows)))
	return activeWindows
}
func SetBorderlessWindow(hwnd uintptr) {
	idxStyle := int32(GWL_STYLE)
	curStyle, _, _ := procGetWindowLong.Call(
		uintptr(hwnd),
		uintptr(idxStyle),
	)
	newStyle := uint32(curStyle)
	newStyle &^= (WS_CAPTION | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX | WS_SYSMENU)
	newStyle |= (WS_POPUP | WS_VISIBLE | WS_CLIPSIBLINGS | WS_CLIPCHILDREN)

	r, _, callErr := procSetWindowLong.Call(
		uintptr(hwnd),
		uintptr(idxStyle),
		uintptr(newStyle),
	)
	if r == 0 && callErr != nil && callErr != syscall.Errno(0) {
		ErrorLog(fmt.Sprintf("SetWindowLongW(GWL_STYLE) failed for hwnd=0x%x: %v", hwnd, callErr))
		return
	}
	idxExStyle := int32(GWL_EXSTYLE)
	curExStyle, _, _ := procGetWindowLong.Call(
		uintptr(hwnd),
		uintptr(idxExStyle),
	)
	newExStyle := uint32(curExStyle)
	newExStyle &^= (WS_EX_DLGMODALFRAME | WS_EX_CLIENTEDGE | WS_EX_STATICEDGE)

	r, _, callErr = procSetWindowLong.Call(
		uintptr(hwnd),
		uintptr(idxExStyle),
		uintptr(newExStyle),
	)
	if r == 0 && callErr != nil && callErr != syscall.Errno(0) {
		ErrorLog(fmt.Sprintf("SetWindowLongW(GWL_EXSTYLE) failed for hwnd=0x%x: %v", hwnd, callErr))
		return
	}

	// Apply the style changes without moving/sizing or changing Z-order
	SetWindowPos(hwnd, 0, 0, 0, 0, 0, SWP_FRAMECHANGED|SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_NOOWNERZORDER|SWP_NOACTIVATE|SWP_SHOWWINDOW)
}
