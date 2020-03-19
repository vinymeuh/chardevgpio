// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Package chardevgpio is a low-level library to the Linux GPIO Character device API.
package chardevgpio

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Chip is a GPIO chips controlling a set of lines.
type Chip struct {
	GPIOChipInfo
	Fd uintptr
}

// Open returns a new Chip for a GPIO character device from its path.
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

// Close releases ressources helded by the Chip.
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

type HandleRequest struct {
	GPIOHandleRequest
}

func (c Chip) RequestOutputLine(line int, consumer string) (HandleRequest, error) {
	hr := HandleRequest{}
	hr.Flags = GPIOHANDLE_REQUEST_OUTPUT
	hr.LineOffsets[0] = uint32(line)
	hr.Lines = 1
	for i, c := range []byte(consumer) {
		if i == 32 {
			break
		}
		hr.Consumer[i] = c
	}
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.Fd, GPIO_GET_LINEHANDLE_IOCTL, uintptr(unsafe.Pointer(&hr.GPIOHandleRequest)))
	if errno != 0 {
		return hr, errno
	}
	return hr, nil
}

func (c Chip) RequestOutputLines(lines []int, consumer string) (HandleRequest, error) {
	hr := HandleRequest{}

	if len(lines) > GPIOHANDLES_MAX {
		return hr, fmt.Errorf("Number of requested lines exceeds GPIOHANDLES_MAX (%d)", GPIOHANDLES_MAX)
	}

	hr.Flags = GPIOHANDLE_REQUEST_OUTPUT

	for _, l := range lines {
		hr.LineOffsets[hr.Lines] = uint32(l)
		hr.Lines++
	}

	for i, c := range []byte(consumer) {
		if i == 32 {
			break
		}
		hr.Consumer[i] = c
	}

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.Fd, GPIO_GET_LINEHANDLE_IOCTL, uintptr(unsafe.Pointer(&hr.GPIOHandleRequest)))
	if errno != 0 {
		return hr, errno
	}
	return hr, nil
}

//func (hr HandleRequest)
