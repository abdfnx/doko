//go:build !windows
// +build !windows

package shared

import (
	"syscall"
	"unsafe"
)

func getTermSize(fd uintptr) (int, int) {
	var sz struct {
		rows uint16
		cols uint16
	}

	_, _, _ = syscall.Syscall(
		syscall.SYS_IOCTL,
		fd, uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&sz)),
	)

	return int(sz.cols), int(sz.rows)
}
