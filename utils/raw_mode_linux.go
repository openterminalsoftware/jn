//go:build linux

package utils

import (
	"os" // Added import
	"syscall"
	"unsafe"
)

func enableRawMode(fd int) (*syscall.Termios, error) {
	var oldState syscall.Termios

	const (
		TCGET = syscall.TCGETS
		TCSET = syscall.TCSETS
	)

	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
		uintptr(TCGET), uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)
	if errno != 0 {
		return nil, errno
	}

	newState := oldState
	newState.Lflag &^= syscall.ICANON | syscall.ECHO
	newState.Iflag &^= syscall.ICRNL
	newState.Oflag &^= syscall.OPOST
	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0

	_, _, errno = syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
		uintptr(TCSET), uintptr(unsafe.Pointer(&newState)), 0, 0, 0)
	if errno != 0 {
		return nil, errno
	}

	return &oldState, nil
}

func disableRawMode(fd int, oldState *syscall.Termios) {
	const TCSET = syscall.TCSETS
	syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
		uintptr(TCSET), uintptr(unsafe.Pointer(oldState)), 0, 0, 0)
}

// ReadChar reads a single character from stdin.
func ReadChar() (rune, int, error) {
	var buf [1]byte
	n, err := os.Stdin.Read(buf[:])
	if err != nil {
		return 0, 0, err
	}
	return rune(buf[0]), n, nil
}