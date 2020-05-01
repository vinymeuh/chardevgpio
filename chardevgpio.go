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

// Name returns the name of the chip.
func (c Chip) Name() string {
	return bytesToString(c.name)

}

// Label returns the label of the chip.
func (c Chip) Label() string {
	return bytesToString(c.label)

}

// Lines returns the number of lines managed by the chip.
func (c Chip) Lines() int {
	return int(c.lines)
}

// Close releases resources helded by the chip.
func (c Chip) Close() error {
	return syscall.Close(int(c.fd))
}

// LineInfo returns informations about the requested line.
func (c Chip) LineInfo(offset int) (LineInfo, error) {
	var li LineInfo
	li.offset = uint32(offset)
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.fd, gpioGetLineInfoIOCTL, uintptr(unsafe.Pointer(&li)))
	if errno != 0 {
		return li, errno
	}
	return li, nil
}

// Offset returns the offset number of the line.
func (li LineInfo) Offset() int {
	return int(li.offset)
}

// Name returns the name of the line.
func (li LineInfo) Name() string {
	return bytesToString(li.name)
}

// Consumer returns the consumer of the line.
func (li LineInfo) Consumer() string {
	return bytesToString(li.consumer)

}

// IsOutput returns true if the line is configured as an output.
func (li LineInfo) IsOutput() bool {
	return li.flags&lineFlagIsOut == lineFlagIsOut
}

// IsInput returns true if the line is configured as an input.
func (li LineInfo) IsInput() bool {
	return !li.IsOutput()
}

// IsActiveLow returns true if the line is configured as active low.
func (li LineInfo) IsActiveLow() bool {
	return li.flags&lineFlagActiveLow == lineFlagActiveLow
}

// IsActiveHigh returns true if the line is configured as active high.
func (li LineInfo) IsActiveHigh() bool {
	return !(li.IsActiveLow())
}

// IsOpenDrain returns true if the line is configured as open drain.
func (li LineInfo) IsOpenDrain() bool {
	return li.flags&lineFlagOpenDrain == lineFlagOpenDrain
}

// IsOpenSource returns true if the line is configured as open source.
func (li LineInfo) IsOpenSource() bool {
	return li.flags&lineFlagOpenSource == lineFlagOpenSource
}

// IsKernel returns true if the line is configured as kernel.
func (li LineInfo) IsKernel() bool {
	return li.flags&lineFlagKernel == lineFlagKernel
}

/*** WIP ***/

// NewHandleRequest prepare a HandleRequest
func NewHandleRequest(offsets []int, flags HandleRequestFlag) *HandleRequest {
	if len(offsets) > handlesMax {
		panic(fmt.Sprintf("Number of requested lines exceeds maximum authorized (%d)", handlesMax))
	}

	hr := &HandleRequest{}
	hr.flags = flags

	for i := range offsets {
		hr.lineOffsets[i] = uint32(offsets[i])
		hr.lines++
	}

	return hr
}

// WithConsumer set the consumer for a prepared HandleRequest.
func (hr *HandleRequest) WithConsumer(consumer string) *HandleRequest {
	hr.consumer = stringToBytes(consumer)
	return hr
}

// WithDefaults set the default values for a prepared HandleRequest.
func (hr *HandleRequest) WithDefaults(defaults []int) *HandleRequest {
	if len(defaults) > handlesMax {
		panic(fmt.Sprintf("Number of default values exceeds maximum authorized (%d)", handlesMax))
	}

	for i := range defaults {
		hr.defaultValues[i] = uint8(defaults[i])
	}
	return hr
}

// RequestLines takes a prepared HandleRequest and returns it ready to work.
func (c Chip) RequestLines(request *HandleRequest) error {
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.fd, gpioGetLineHandleIOCTL, uintptr(unsafe.Pointer(request)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Write writes values to lines handled by the HandleRequest.
func (hr *HandleRequest) Write(values []int) error {
	out := Data{}
	for i := range values {
		if i > handlesMax-1 {
			break
		}
		out.Values[i] = uint8(values[i])
	}

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(hr.fd), gpioHandleSetLineValuesIOCTL, uintptr(unsafe.Pointer(&out)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Write0 writes value to the first line handled by the HandleRequest.
func (hr *HandleRequest) Write0(value int) error {
	out := Data{}
	out.Values[0] = uint8(value)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(hr.fd), gpioHandleSetLineValuesIOCTL, uintptr(unsafe.Pointer(&out)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Close releases resources helded by the HandleRequest.
func (hr HandleRequest) Close() error {
	return syscall.Close(int(hr.fd))
}

/*** WIP ***/

// bytesToString is a helper function to convert raw string as stored in Linux structure to Go string.
func bytesToString(B [32]byte) string {
	n := bytes.IndexByte(B[:], 0)
	if n == -1 {
		return string(B[:])
	}
	return string(B[:n])
}

// stringToBytes is a helper function to convert Go string to string as stored in Linux structure.
// Used to set GPIOHandleRequest.Consumer
func stringToBytes(s string) [32]byte {
	var b [32]byte
	for i, c := range []byte(s) {
		if i == 32 {
			break
		}
		b[i] = c
	}
	return b
}
