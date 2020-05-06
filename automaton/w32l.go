package automaton

// contains functions that aren't in JamesHovious/w32
// TODO: fork JamesHovious/w32 and add the following functions

import (
	"fmt"
	"syscall"

	"github.com/JamesHovious/w32"
)

var (
	modUser32         = syscall.NewLazyDLL("user32.dll")
	modGdi32          = syscall.NewLazyDLL("gdi32.dll")
	procPrintWindow   = modUser32.NewProc("PrintWindow")
	procGetDIBits     = modGdi32.NewProc("GetDIBits")
	procGetTopWindow  = modUser32.NewProc("GetTopWindow")
	procGetNextWindow = modUser32.NewProc("GetWindow")
)

type WindowsError struct {
	Msg string
}

func (w *WindowsError) Error() string {
	return fmt.Sprintf("WindowsError: %s", w.Msg)
}

func PrintWindow(hwnd w32.HWND, hdc w32.HDC) bool {
	ret, _, err := procPrintWindow.Call(
		uintptr(hwnd), uintptr(hdc), 1)
	fmt.Printf("PrintWindow: %v %v\n", ret, err)
	return ret != 0
}

func GetDIBits(hdc w32.HDC, hBmp w32.HBITMAP, start uint32,
	cLines int32, buf uintptr, bmpInfo uintptr, usage uint32) int {
	ret, _, err := procGetDIBits.Call(
		uintptr(hdc), uintptr(hBmp), uintptr(start), uintptr(cLines), buf, bmpInfo, uintptr(usage))
	fmt.Printf("GetDIBits: %v %v\n", ret, err)
	return int(ret)
}

func GetTopWindow(hwnd w32.HWND) (w32.HWND, error) {
	ret, _, _ := procGetTopWindow.Call(uintptr(hwnd))

	if ret == 0 {
		return 0, &WindowsError{"Specified window has no child windows."}
	}

	return w32.HWND(ret), nil
}

func GetNextWindow(hwnd w32.HWND) w32.HWND {
	ret, _, _ := procGetNextWindow.Call(uintptr(hwnd), 2)
	return w32.HWND(ret)
}

func GetAllTopWindows() map[string]w32.HWND {
	hwnd, err := GetTopWindow(w32.HWND(0))
	if err != nil {
		panic(err)
	}

	m := make(map[string]w32.HWND)
	for hwnd != w32.HWND(0) {
		windowText := w32.GetWindowText(hwnd)

		if windowText != "" {
			m[windowText] = hwnd
		}

		hwnd = GetNextWindow(hwnd)
	}

	return m
}
