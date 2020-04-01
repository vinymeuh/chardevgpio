// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// +build linux

// Package chardevgpio is a low-level library to the Linux GPIO Character device API.
package chardevgpio

import (
	"bytes"
	"encoding/binary"
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

// LineDirection is used to indicate if the requested line will be used as an input or an output.
type LineDirection uint32

const (
	// LineIn setup a line as an input.
	LineIn LineDirection = GPIOHANDLE_REQUEST_INPUT
	// LineOut setup a line as an output.
	LineOut = GPIOHANDLE_REQUEST_OUTPUT
)

// DataLine represents a single line to be used to send or receive data.
type DataLine struct {
	GPIOHandleRequest
}

// RequestOutputLine requests to the chip a single DataLine to send data.
func (c Chip) RequestOutputLine(line int, value int, consumer string) (DataLine, error) {
	l := DataLine{}
	l.Flags = GPIOHANDLE_REQUEST_OUTPUT
	l.LineOffsets[0] = uint32(line)
	l.DefaultValues[0] = uint8(value)
	l.Lines = 1
	l.Consumer = consumerFromString(consumer)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.Fd, GPIO_GET_LINEHANDLE_IOCTL, uintptr(unsafe.Pointer(&l.GPIOHandleRequest)))
	if errno != 0 {
		return l, errno
	}
	return l, nil
}

// RequestInputLine requests to the chip a single DataLine to receive data.
func (c Chip) RequestInputLine(line int, consumer string) (DataLine, error) {
	l := DataLine{}
	l.Flags = GPIOHANDLE_REQUEST_INPUT
	l.LineOffsets[0] = uint32(line)
	l.Lines = 1
	l.Consumer = consumerFromString(consumer)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.Fd, GPIO_GET_LINEHANDLE_IOCTL, uintptr(unsafe.Pointer(&l.GPIOHandleRequest)))
	if errno != 0 {
		return l, errno
	}
	return l, nil
}

// Close releases resources helded by the DataLine.
func (l DataLine) Close() error {
	return syscall.Close(l.Fd)
}

// SetValue writes value to a DataLine.
func (l DataLine) SetValue(value int) error {
	hd := GPIOHandleData{}
	hd.Values[0] = uint8(value)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(l.Fd), GPIOHANDLE_SET_LINE_VALUES_IOCTL, uintptr(unsafe.Pointer(&hd)))
	if errno != 0 {
		return errno
	}
	return nil
}

// DataLines represents a set of lines to be used to send or receive data.
type DataLines struct {
	GPIOHandleRequest
}

// RequestOutputLines requests to the chip a DataLines to receive data.
func (c Chip) RequestOutputLines(lines []int, values []int, consumer string) (DataLines, error) {
	L := DataLines{}
	if len(lines) > GPIOHANDLES_MAX {
		return L, fmt.Errorf("Number of requested lines exceeds GPIOHANDLES_MAX (%d)", GPIOHANDLES_MAX)
	}
	if len(values) < len(lines) {
		return L, fmt.Errorf("Not enough values to initialize lines")
	}

	L.Flags = GPIOHANDLE_REQUEST_OUTPUT
	for i := range lines {
		L.LineOffsets[i] = uint32(lines[i])
		L.DefaultValues[i] = uint8(values[i])
		L.Lines++
	}
	L.Consumer = consumerFromString(consumer)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.Fd, GPIO_GET_LINEHANDLE_IOCTL, uintptr(unsafe.Pointer(&L.GPIOHandleRequest)))
	if errno != 0 {
		return L, errno
	}
	return L, nil
}

// RequestInputLines requests to the chip a DataLines to send data.
func (c Chip) RequestInputLines(lines []int, consumer string) (DataLines, error) {
	L := DataLines{}
	if len(lines) > GPIOHANDLES_MAX {
		return L, fmt.Errorf("Number of requested lines exceeds GPIOHANDLES_MAX (%d)", GPIOHANDLES_MAX)
	}

	L.Flags = GPIOHANDLE_REQUEST_INPUT
	for i := range lines {
		L.LineOffsets[i] = uint32(lines[i])
		L.Lines++
	}
	L.Consumer = consumerFromString(consumer)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, c.Fd, GPIO_GET_LINEHANDLE_IOCTL, uintptr(unsafe.Pointer(&L.GPIOHandleRequest)))
	if errno != 0 {
		return L, errno
	}
	return L, nil
}

// SetValues writes value to a DataLines.
func (L DataLines) SetValues(values []int) error {
	hd := GPIOHandleData{}
	for i := range values {
		if i > GPIOHANDLES_MAX-1 {
			break
		}
		hd.Values[i] = uint8(values[i])
	}

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(L.Fd), GPIOHANDLE_SET_LINE_VALUES_IOCTL, uintptr(unsafe.Pointer(&hd)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Close releases resources helded by the DataLines.
func (L DataLines) Close() error {
	return syscall.Close(L.Fd)
}

// EventLineType defines the kind of event to wait on an event line.
type EventLineType uint32

// EventLineType
const (
	RisingEdge   EventLineType = GPIOEVENT_REQUEST_RISING_EDGE
	FaillingEdge               = GPIOEVENT_REQUEST_FALLING_EDGE
	BothEdge                   = GPIOEVENT_REQUEST_BOTH_EDGES
)

// EventLineWatcher waits for events on a set of event lines.
type EventLineWatcher struct {
	epfd int
	efds []int
}

// NewEventLineWatcher initializes a new EventLineWatcher.
func NewEventLineWatcher() (EventLineWatcher, error) {
	fd, err := unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	return EventLineWatcher{epfd: fd}, err
}

// AddEvent adds a new event line to watch to the EventLineWatcher.
func (elw *EventLineWatcher) AddEvent(chip Chip, line int, consumer string, eventType EventLineType) error {
	el := GPIOEventRequest{}
	el.Consumer = consumerFromString(consumer)
	el.LineOffset = uint32(line)
	el.HandleFlags = GPIOHANDLE_REQUEST_INPUT
	el.EventFlags = uint32(eventType)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, chip.Fd, GPIO_GET_LINEEVENT_IOCTL, uintptr(unsafe.Pointer(&el)))
	if errno != 0 {
		return errno
	}
	// an application that employs the EPOLLET flag should use nonblocking file descriptors (man epoll)
	unix.SetNonblock(el.Fd, true)

	// add the EventLine fd to the epoll instance
	var epEvent unix.EpollEvent
	epEvent.Events = unix.EPOLLIN | unix.EPOLLET
	epEvent.Fd = int32(el.Fd)
	if err := unix.EpollCtl(elw.epfd, unix.EPOLL_CTL_ADD, el.Fd, &epEvent); err != nil {
		return err
	}
	elw.efds = append(elw.efds, el.Fd)

	return nil
}

// EventData is a type alias for GPIOEventData
type EventData = GPIOEventData

// EventLineHandler is the type of the function called for each EventData retrieved by Walk or WalkForEver.
type EventLineHandler func(evd EventData)

// Wait waits for first occurrence of an event on the event lines.
// Note that for one event, more than one EventData can be retrieved on the event line.
func (elw *EventLineWatcher) Wait(handler EventLineHandler) error {
	var events [1]unix.EpollEvent
	for {
		_, err := unix.EpollWait(elw.epfd, events[:], -1)
		if err == nil {
			break
		} else {
			if err == unix.EINTR {
				continue
			}
			return err
		}
	}

	ev := events[0]
	if ev.Events&unix.EPOLLIN != 0 {
		evds, err := readEventsData(int(ev.Fd))
		if err != nil {
			return err
		}
		for _, evd := range evds {
			handler(evd)
		}
	}

	return nil
}

// WaitForEver waits indefinitely for events on then event lines.
func (elw *EventLineWatcher) WaitForEver(handler EventLineHandler) error {
	events := make([]unix.EpollEvent, len(elw.efds))
	for {
		nevents, err := unix.EpollWait(elw.epfd, events[:], -1)
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

// Close releases resources helded by the EventLineWatcher.
func (elw *EventLineWatcher) Close() error {
	err := unix.Close(elw.epfd)
	for _, fd := range elw.efds {
		unix.Close(fd) // TODO: concatenate errors
	}
	return err
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

// helper that retrieves all event data that can be retrieved on a event line.
// The event line is fully drained when read receives EAGAIN.
func readEventsData(fd int) ([]GPIOEventData, error) {
	// How to know that buffer size must be 16 ?
	// GPIOEventData = uint64 + uint32 = 8 + 4 = 12 !?
	const BUFFER_SIZE = 16

	var evds []GPIOEventData
	var evd GPIOEventData
	var buffer = make([]byte, BUFFER_SIZE)
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
