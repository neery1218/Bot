package automaton

// contains functions that aren't in JamesHovious/w32
// TODO: fork JamesHovious/w32 and add the following functions

import (
	"fmt"
	"github.com/JamesHovious/w32"
	"syscall"
)

var (
	modUser32       = syscall.NewLazyDLL("user32.dll")
	modGdi32        = syscall.NewLazyDLL("gdi32.dll")
	procPrintWindow = modUser32.NewProc("PrintWindow")
	procGetDIBits   = modGdi32.NewProc("GetDIBits")
)

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
