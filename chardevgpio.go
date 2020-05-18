// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// +build linux

// Package chardevgpio is a library to the Linux GPIO Character device API.
package chardevgpio

import (
	"bytes"
	"encoding/binary"
	"errors"
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
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.fd, ioctlGetChipInfo, uintptr(unsafe.Pointer(&c.ChipInfo)))
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
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.fd, ioctlGetLineInfo, uintptr(unsafe.Pointer(&li)))
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

// IsBiasPullUp returns true if the line is configured as bias pull up.
func (li LineInfo) IsBiasPullUp() bool {
	return li.flags&lineFlagBiasPullUp == lineFlagBiasPullUp
}

// IsBiasPullDown returns true if the line is configured as bias pull down.
func (li LineInfo) IsBiasPullDown() bool {
	return li.flags&lineFlagBiasPullDown == lineFlagBiasPullDown
}

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
// Note that setting default values on a InputLine is a nonsense but no error are returned.
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
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.fd, ioctlGetLineHandle, uintptr(unsafe.Pointer(request)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Reads return values read from the lines handled by the HandleRequest.
// The second return parameter contains all values returned as an array.
// The first one is the first element of this array, useful when dealing with 1 line HandleRequest.
func (hr *HandleRequest) Read() (int, []int, error) {
	if hr.flags&lineFlagIsOut == lineFlagIsOut {
		return 0, []int{}, ErrOperationNotPermitted
	}

	in := handleData{}
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(hr.fd), ioctlHandleGetLineValues, uintptr(unsafe.Pointer(&in)))
	if errno != 0 {
		return 0, []int{}, errno
	}

	switch len(in.values) {
	case 1:
		return int(in.values[0]), []int{int(in.values[0])}, nil
	default:
		valueN := make([]int, len(in.values))
		for i := range in.values {
			valueN[i] = int(in.values[i])
		}
		return valueN[0], valueN, nil
	}
}

// Write writes values to the lines handled by the HandleRequest.
// If there is more values ​​supplied than lines managed by the HandleRequest, excess values ​​are silently ignored.
func (hr *HandleRequest) Write(value0 int, valueN ...int) error {
	if !(hr.flags&lineFlagIsOut == lineFlagIsOut) {
		return ErrOperationNotPermitted
	}

	out := handleData{}
	out.values[0] = uint8(value0)
	for i := range valueN {
		if i >= int(hr.lines)-1 {
			break
		}
		out.values[i+1] = uint8(valueN[i])
	}

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(hr.fd), ioctlHandleSetLineValues, uintptr(unsafe.Pointer(&out)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Close releases resources helded by the HandleRequest.
func (hr HandleRequest) Close() error {
	return syscall.Close(int(hr.fd))
}

// IsRising returns true for event on a rising edge.
func (e Event) IsRising() bool {
	return e.ID == eventRisingEdge
}

// IsFalling returns true for event on a falling edge.
func (e Event) IsFalling() bool {
	return e.ID == eventFallingEdge
}

// LineWatcher is a receiver of events for a set of event lines.
type LineWatcher struct {
	epfd int
	efds []int
}

// NewLineWatcher initializes a new LineWatcher.
func NewLineWatcher() (LineWatcher, error) {
	fd, err := unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	return LineWatcher{epfd: fd}, err
}

// Close releases resources helded by the LineWatcher.
func (lw *LineWatcher) Close() error {
	err := unix.Close(lw.epfd)
	for _, fd := range lw.efds {
		unix.Close(fd) // TODO: concatenate errors
	}
	return err
}

// Add adds a new line to watch to the LineWatcher.
func (lw *LineWatcher) Add(chip Chip, line int, flags EventRequestFlags, consumer string) error {
	el := EventLine{
		lineOffset:  uint32(line),
		handleFlags: HandleRequestInput,
		eventFlags:  uint32(flags),
		consumer:    stringToBytes(consumer),
	}

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, chip.fd, ioctlGetLineEvent, uintptr(unsafe.Pointer(&el)))
	if errno != 0 {
		return errno
	}
	// an application that employs the EPOLLET flag should use nonblocking file descriptors (man epoll)
	unix.SetNonblock(int(el.fd), true)

	// add the EventLine fd to the epoll instance
	var epEvent unix.EpollEvent
	epEvent.Events = unix.EPOLLIN | unix.EPOLLET
	epEvent.Fd = int32(el.fd)
	if err := unix.EpollCtl(lw.epfd, unix.EPOLL_CTL_ADD, int(el.fd), &epEvent); err != nil {
		return err
	}
	lw.efds = append(lw.efds, int(el.fd))

	return nil
}

// Wait waits for first occurrence of an event on one of the event lines.
func (lw *LineWatcher) Wait() (Event, error) {
	var events [1]unix.EpollEvent
	for {
		_, err := unix.EpollWait(lw.epfd, events[:], -1)
		if err != nil {
			if err == unix.EINTR {
				continue
			}
			return Event{}, err
		}

		ev := events[0]
		if ev.Events&unix.EPOLLIN != 0 {
			evds, err := readEventsData(int(ev.Fd))
			if err != nil {
				return Event{}, err
			}
			return evds[0], nil
		}
	}
}

// EventHandlerFunc is the type of the function called for each event retrieved by WaitForEver.
type EventHandlerFunc func(evd Event)

// WaitForEver waits indefinitely for events on the event lines.
// Note that for one event, more than one EventData can be retrieved on the event line.
func (lw *LineWatcher) WaitForEver(handler EventHandlerFunc) error {
	events := make([]unix.EpollEvent, len(lw.efds))
	for {
		nevents, err := unix.EpollWait(lw.epfd, events[:], -1)
		if err != nil {
			if err == unix.EINTR {
				continue
			}
			return err
		}

		for i := 0; i < nevents; i++ {
			ev := events[i]
			if ev.Events&unix.EPOLLIN != 0 {
				evds, err := readEventsData(int(ev.Fd))
				if err != nil {
					return err
				}
				for _, evd := range evds {
					handler(evd)
				}
			}
		}
	}
}

// readEventsData that retrieves all event data that can be retrieved on a event line.
// The event line is fully drained when read receives EAGAIN.
func readEventsData(fd int) ([]Event, error) {
	const BufferSize = 16 // How to know that buffer size must be 16, GPIOEventData = uint64 + uint32 = 8 + 4 = 12 ?

	var evds []Event
	var evd Event
	var buffer = make([]byte, BufferSize)
	for {
		_, err := unix.Read(fd, buffer)
		if err != nil {
			if err == unix.EAGAIN {
				return evds, nil
			}
			return evds, err
		}

		err = binary.Read(bytes.NewReader(buffer), binary.LittleEndian, &evd)
		if err != nil {
			return evds, err
		}
		evds = append(evds, evd)
	}
}

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

// ErrOperationNotPermitted is returned when trying to read on an output line or to write on a input line.
var ErrOperationNotPermitted = errors.New("operation not permitted")
