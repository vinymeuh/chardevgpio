// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// +build linux

package chardevgpio

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

/** DataLine structure is defined in linux.go **/

// SetValue writes value to a DataLine.
func (l DataLine) SetValue(value int) error {
	d := Data{}
	d.Values[0] = uint8(value)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(l.Fd), gpioHandleSetLineValuesIOCTL, uintptr(unsafe.Pointer(&d)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Close releases resources helded by the DataLine.
func (l DataLine) Close() error {
	return syscall.Close(int(l.Fd))
}

// DataLines represents a set of lines to be used to send or receive data.
type DataLines struct {
	DataLine
}

// SetValues writes value to a DataLines.
func (L DataLines) SetValues(values []int) error {
	D := Data{}
	for i := range values {
		if i > gpioHandlesMax-1 {
			break
		}
		D.Values[i] = uint8(values[i])
	}

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(L.Fd), gpioHandleSetLineValuesIOCTL, uintptr(unsafe.Pointer(&D)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Close releases resources helded by the DataLines.
func (L DataLines) Close() error {
	return syscall.Close(int(L.Fd))
}
