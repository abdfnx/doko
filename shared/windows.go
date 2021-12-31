//go:build windows
// +build windows

package shared

import (
	"syscall"
	"unsafe"
)

func getTermSize(fd uintptr) (int, int) {
	var csbi consoleScreenBufferInfo

	r1, _, _ := procGetConsoleScreenBufferInfo.Call(fd, uintptr(unsafe.Pointer(&csbi)))

	if r1 == 0 {
		return 80, 25
	}

	return int(csbi.window.right - csbi.window.left + 1), int(csbi.window.bottom - csbi.window.top + 1)
}

func IsTermWindowSizeBiggerThanZero() bool {
	h, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)

	if err != nil {
		return true
	}

	termw, termh := getTermSize(uintptr(h))

	if termw > 0 && termh > 0 {
		return true
	}

	return false
}
