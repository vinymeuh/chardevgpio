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

// Chip is a GPIO chip controlling a set of lines.
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

// Close releases resources helded by the Chip.
func (c Chip) Close() error {
	return syscall.Close(int(c.Fd))
}

// LineInfo contains informations about a line.
type LineInfo struct {
	GPIOLineInfo
}

// LineInfo returns informations about the requested line.
func (c Chip) LineInfo(line int) (LineInfo, error) {
	var l LineInfo
	l.GPIOLineInfo.Offset = uint32(line)
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.Fd, GPIO_GET_LINEINFO_IOCTL, uintptr(unsafe.Pointer(&l.GPIOLineInfo)))
	if errno != 0 {
		return l, errno
	}
	return l, nil
}

// IsOutput returns true if the line is configured as an output.
func (li LineInfo) IsOutput() bool {
	return li.GPIOLineInfo.Flags&GPIOLINE_FLAG_IS_OUT == GPIOLINE_FLAG_IS_OUT
}

// IsInput returns true if the line is configured as an input.
func (li LineInfo) IsInput() bool {
	return !(li.GPIOLineInfo.Flags&GPIOLINE_FLAG_IS_OUT == GPIOLINE_FLAG_IS_OUT)
}

// IsActiveLow returns true if the line is configured as active low.
func (li LineInfo) IsActiveLow() bool {
	return li.GPIOLineInfo.Flags&GPIOLINE_FLAG_ACTIVE_LOW == GPIOLINE_FLAG_ACTIVE_LOW
}

// IsActiveHigh returns true if the line is configured as active high.
func (li LineInfo) IsActiveHigh() bool {
	return !(li.GPIOLineInfo.Flags&GPIOLINE_FLAG_ACTIVE_LOW == GPIOLINE_FLAG_ACTIVE_LOW)
}

// IsOpenDrain returns true if the line is configured as open drain.
func (li LineInfo) IsOpenDrain() bool {
	return li.GPIOLineInfo.Flags&GPIOLINE_FLAG_OPEN_DRAIN == GPIOLINE_FLAG_OPEN_DRAIN
}

// IsOpenSource returns true if the line is configured as open source.
func (li LineInfo) IsOpenSource() bool {
	return li.GPIOLineInfo.Flags&GPIOLINE_FLAG_OPEN_SOURCE == GPIOLINE_FLAG_OPEN_SOURCE
}

// IsKernel returns true if the line is configured as kernel.
func (li LineInfo) IsKernel() bool {
	return li.GPIOLineInfo.Flags&GPIOLINE_FLAG_KERNEL == GPIOLINE_FLAG_KERNEL
}

// Line represents a single requested line.
type Line struct {
	GPIOHandleRequest
}

func (c Chip) RequestOutputLine(line int, consumer string) (Line, error) {
	l := Line{}
	l.Consumer = consumerFromString(consumer)
	l.Flags = GPIOHANDLE_REQUEST_OUTPUT
	l.LineOffsets[0] = uint32(line)
	l.Lines = 1

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.Fd, GPIO_GET_LINEHANDLE_IOCTL, uintptr(unsafe.Pointer(&l.GPIOHandleRequest)))
	if errno != 0 {
		return l, errno
	}
	return l, nil
}

// Close releases resources helded by the Line.
func (l Line) Close() error {
	return syscall.Close(l.Fd)
}

func (l Line) SetValue(value int) error {
	hd := GPIOHandleData{}
	hd.Values[0] = uint8(value)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(l.Fd), GPIOHANDLE_SET_LINE_VALUES_IOCTL, uintptr(unsafe.Pointer(&hd)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Lines represents a set of requested lines.
type Lines struct {
	GPIOHandleRequest
}

func (c Chip) RequestOutputLines(lines []int, consumer string) (Lines, error) {
	L := Lines{}

	if len(lines) > GPIOHANDLES_MAX {
		return L, fmt.Errorf("Number of requested lines exceeds GPIOHANDLES_MAX (%d)", GPIOHANDLES_MAX)
	}

	L.Consumer = consumerFromString(consumer)
	L.Flags = GPIOHANDLE_REQUEST_OUTPUT
	for _, l := range lines {
		L.LineOffsets[L.Lines] = uint32(l)
		L.Lines++
	}

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.Fd, GPIO_GET_LINEHANDLE_IOCTL, uintptr(unsafe.Pointer(&L.GPIOHandleRequest)))
	if errno != 0 {
		return L, errno
	}
	return L, nil
}

// Close releases resources helded by the Lines.
func (L Lines) Close() error {
	return syscall.Close(L.Fd)
}

// helper that convert a string to an array of 32 bytes
// Used to set GPIOHandleRequest.Consumer
func consumerFromString(consumer string) [32]byte {
	var b [32]byte
	for i, c := range []byte(consumer) {
		if i == 32 {
			break
		}
		b[i] = c
	}
	return b
}
