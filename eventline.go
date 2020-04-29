// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// +build linux

package chardevgpio

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"golang.org/x/sys/unix"
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
	el := EventLine{}
	el.Consumer = consumerFromString(consumer)
	el.LineOffset = uint32(line)
	el.HandleFlags = gpioHandleRequestInput
	el.EventFlags = uint32(eventType)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, chip.fd, gpioGetLineEventIOCTL, uintptr(unsafe.Pointer(&el)))
	if errno != 0 {
		return errno
	}
	// an application that employs the EPOLLET flag should use nonblocking file descriptors (man epoll)
	unix.SetNonblock(int(el.Fd), true)

	// add the EventLine fd to the epoll instance
	var epEvent unix.EpollEvent
	epEvent.Events = unix.EPOLLIN | unix.EPOLLET
	epEvent.Fd = int32(el.Fd)
	if err := unix.EpollCtl(elw.epfd, unix.EPOLL_CTL_ADD, int(el.Fd), &epEvent); err != nil {
		return err
	}
	elw.efds = append(elw.efds, int(el.Fd))

	return nil
}

// EventLineHandler is the type of the function called for each Event retrieved by Walk or WalkForEver.
type EventLineHandler func(evd Event)

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
func readEventsData(fd int) ([]Event, error) {
	// How to know that buffer size must be 16 ?
	// GPIOEventData = uint64 + uint32 = 8 + 4 = 12 !?
	const BufferSize = 16

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
