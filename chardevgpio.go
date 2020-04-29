// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// +build linux

// Package chardevgpio is a library to the Linux GPIO Character device API.
package chardevgpio

import (
	"bytes"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Chip is a GPIO chip controlling a set of lines.
type Chip struct {
	ChipInfo
	fd uintptr
}

// NewChip returns a Chip for a GPIO character device from its path.
func NewChip(path string) (Chip, error) {
	f, err := os.Open(path)
	if err != nil {
		return Chip{}, err
	}

	var c Chip
	c.fd = f.Fd()
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.fd, gpioGetChipInfoIOCTL, uintptr(unsafe.Pointer(&c.ChipInfo)))
	if errno != 0 {
		return c, errno
	}
	return c, nil
}

// Name returns the name of the Chip.
func (c Chip) Name() string {
	n := bytes.IndexByte(c.name[:], 0)
	if n == -1 {
		return string(c.name[:])
	}
	return string(c.name[:n])
}

// Label returns the label of the Chip.
func (c Chip) Label() string {
	n := bytes.IndexByte(c.label[:], 0)
	if n == -1 {
		return string(c.label[:])
	}
	return string(c.label[:n])
}

// Lines returns the number of lines managed by the Chip.
func (c Chip) Lines() int {
	return int(c.lines)
}

// Close releases resources helded by the Chip.
func (c Chip) Close() error {
	return syscall.Close(int(c.fd))
}

// LineInfo returns informations about the requested line.
func (c Chip) LineInfo(line int) (LineInfo, error) {
	var li LineInfo
	li.Offset = uint32(line)
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.fd, gpioGetLineInfoIOCTL, uintptr(unsafe.Pointer(&li)))
	if errno != 0 {
		return li, errno
	}
	return li, nil
}

// IsOutput returns true if the line is configured as an output.
func (li LineInfo) IsOutput() bool {
	return li.Flags&gpioLineFlagIsOut == gpioLineFlagIsOut
}

// IsInput returns true if the line is configured as an input.
func (li LineInfo) IsInput() bool {
	return !li.IsOutput()
}

// IsActiveLow returns true if the line is configured as active low.
func (li LineInfo) IsActiveLow() bool {
	return li.Flags&gpioLineFlagActiveLow == gpioLineFlagActiveLow
}

// IsActiveHigh returns true if the line is configured as active high.
func (li LineInfo) IsActiveHigh() bool {
	return !(li.IsActiveLow())
}

// IsOpenDrain returns true if the line is configured as open drain.
func (li LineInfo) IsOpenDrain() bool {
	return li.Flags&gpioLineFlagOpenDrain == gpioLineFlagOpenDrain
}

// IsOpenSource returns true if the line is configured as open source.
func (li LineInfo) IsOpenSource() bool {
	return li.Flags&gpioLineFlagOpenSource == gpioLineFlagOpenSource
}

// IsKernel returns true if the line is configured as kernel.
func (li LineInfo) IsKernel() bool {
	return li.Flags&gpioLineFlagKernel == gpioLineFlagKernel
}

// RequestOutputLine requests to the chip a single DataLine to send data.
func (c Chip) RequestOutputLine(line int, value int, consumer string) (DataLine, error) {
	l := DataLine{}
	l.Flags = gpioHandleRequestOutput
	l.LineOffsets[0] = uint32(line)
	l.DefaultValues[0] = uint8(value)
	l.Lines = 1
	l.Consumer = consumerFromString(consumer)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.fd, gpioGetLineHandleIOCTL, uintptr(unsafe.Pointer(&l)))
	if errno != 0 {
		return l, errno
	}
	return l, nil
}

// RequestInputLine requests to the chip a single DataLine to receive data.
func (c Chip) RequestInputLine(line int, consumer string) (DataLine, error) {
	l := DataLine{}
	l.Flags = gpioHandleRequestInput
	l.LineOffsets[0] = uint32(line)
	l.Lines = 1
	l.Consumer = consumerFromString(consumer)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.fd, gpioGetLineHandleIOCTL, uintptr(unsafe.Pointer(&l)))
	if errno != 0 {
		return l, errno
	}
	return l, nil
}

// RequestOutputLines requests to the chip a DataLines to receive data.
func (c Chip) RequestOutputLines(lines []int, values []int, consumer string) (DataLines, error) {
	L := DataLines{}
	if len(lines) > gpioHandlesMax {
		return L, fmt.Errorf("Number of requested lines exceeds maximum authorized (%d)", gpioHandlesMax)
	}
	if len(values) < len(lines) {
		return L, fmt.Errorf("Not enough values to initialize lines")
	}

	L.Flags = gpioHandleRequestOutput
	for i := range lines {
		L.LineOffsets[i] = uint32(lines[i])
		L.DefaultValues[i] = uint8(values[i])
		L.Lines++
	}
	L.Consumer = consumerFromString(consumer)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.fd, gpioGetLineHandleIOCTL, uintptr(unsafe.Pointer(&L)))
	if errno != 0 {
		return L, errno
	}
	return L, nil
}

// RequestInputLines requests to the chip a DataLines to send data.
func (c Chip) RequestInputLines(lines []int, consumer string) (DataLines, error) {
	L := DataLines{}
	if len(lines) > gpioHandlesMax {
		return L, fmt.Errorf("Number of requested lines exceeds maximum authorized (%d)", gpioHandlesMax)
	}

	L.Flags = gpioHandleRequestInput
	for i := range lines {
		L.LineOffsets[i] = uint32(lines[i])
		L.Lines++
	}
	L.Consumer = consumerFromString(consumer)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.fd, gpioGetLineHandleIOCTL, uintptr(unsafe.Pointer(&L)))
	if errno != 0 {
		return L, errno
	}
	return L, nil
}
