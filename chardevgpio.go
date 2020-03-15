// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Package chardevgpio is a low-level library to the Linux GPIO Character device API.
package chardevgpio

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

type Chip struct {
	GPIOChipInfo
	Fd uintptr
}

func Open(path string) (Chip, error) {
	f, err := os.Open(path)
	if err != nil {
		return Chip{}, err
	}

	var c Chip
	c.Fd = f.Fd()
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.Fd, GPIO_GET_CHIPINFO_IOCTL, uintptr(unsafe.Pointer(&c.GPIOChipInfo)))
	if errno != 0 {
		return c, errno
	}
	return c, nil
}

func (c Chip) Close() error {
	return syscall.Close(int(c.Fd))
}

type Line struct {
	GPIOLineInfo
}

func (c Chip) GetLine(offset int) (Line, error) {
	var l Line
	l.GPIOLineInfo.Offset = uint32(offset)
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.Fd, GPIO_GET_LINEINFO_IOCTL, uintptr(unsafe.Pointer(&l.GPIOLineInfo)))
	if errno != 0 {
		return l, errno
	}
	return l, nil
}

func (l Line) IsOut() bool {
	return l.GPIOLineInfo.Flags&GPIOLINE_FLAG_IS_OUT == GPIOLINE_FLAG_IS_OUT
}

func (l Line) IsIn() bool {
	return !(l.GPIOLineInfo.Flags&GPIOLINE_FLAG_IS_OUT == GPIOLINE_FLAG_IS_OUT)
}

func (l Line) IsActiveLow() bool {
	return l.GPIOLineInfo.Flags&GPIOLINE_FLAG_ACTIVE_LOW == GPIOLINE_FLAG_ACTIVE_LOW
}

func (l Line) IsActiveHigh() bool {
	return !(l.GPIOLineInfo.Flags&GPIOLINE_FLAG_ACTIVE_LOW == GPIOLINE_FLAG_ACTIVE_LOW)
}

func (l Line) IsOpenDrain() bool {
	return l.GPIOLineInfo.Flags&GPIOLINE_FLAG_OPEN_DRAIN == GPIOLINE_FLAG_OPEN_DRAIN
}

func (l Line) IsOpenSource() bool {
	return l.GPIOLineInfo.Flags&GPIOLINE_FLAG_OPEN_SOURCE == GPIOLINE_FLAG_OPEN_SOURCE
}

func (l Line) IsKernel() bool {
	return l.GPIOLineInfo.Flags&GPIOLINE_FLAG_KERNEL == GPIOLINE_FLAG_KERNEL
}
