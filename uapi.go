// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// +build linux

package chardevgpio

import "unsafe"

// Code in this file mimics directly the Linux kernel code

/**
 * ioctl constants from uapi/asm-generic/ioctl.h
 * For reference see https://elixir.bootlin.com/linux/v5.5.9/source/include/uapi/asm-generic/ioctl.h
 */
const (
	iocNRBits   = 8
	iocTypeBits = 8

	iocSizeBits = 14
	iocDirBits  = 2

	iocNRShift   = 0
	iocTypeShift = iocNRShift + iocNRBits
	iocSizeShift = iocTypeShift + iocTypeBits
	iocDirShift  = iocSizeShift + iocSizeBits

	iocRead  = 2
	iocWrite = 1
)

/*
 * gpio code from uapi/linux/gpio.h
 * For reference see https://elixir.bootlin.com/linux/v5.5.9/source/include/uapi/linux/gpio.h
 */

// ChipInfo contains informations about a GPIO chip.
type ChipInfo struct {
	name  [32]byte
	label [32]byte
	lines uint32
}

// Informational flags
const (
	lineFlagKernel       = 1 << 0
	lineFlagIsOut        = 1 << 1
	lineFlagActiveLow    = 1 << 2
	lineFlagOpenDrain    = 1 << 3
	lineFlagOpenSource   = 1 << 4
	lineFlagBiasPullUp   = 1 << 5
	lineFlagBiasPullDown = 1 << 6
	lineFlagDisable      = 1 << 7
)

// LineInfo contains informations about a GPIO line.
type LineInfo struct {
	offset   uint32
	flags    uint32
	name     [32]byte
	consumer [32]byte
}

// handlesMax limits maximum number of handles that can be requested in a GPIOHandleRequest
const handlesMax = 64

// HandleRequestFlag is the type of request flags.
type HandleRequestFlag uint32

// HandleRequest request flags.
const (
	HandleRequestInput        HandleRequestFlag = 1 << 0
	HandleRequestOutput                         = 1 << 1
	HandleRequestActiveLow                      = 1 << 2
	HandleRequestOpenDrain                      = 1 << 3
	HandleRequestOpenSource                     = 1 << 4
	HandleRequestBiasPullUp                     = 1 << 5
	HandleRequestBiasPullDown                   = 1 << 6
	HandleRequestBiasDisable                    = 1 << 7
)

// HandleRequest represents at first a query to be sent to a chip to get control on a set of lines.
// After be returned by the chip, it must be used to send or received data to lines.
type HandleRequest struct {
	lineOffsets   [handlesMax]uint32
	flags         HandleRequestFlag
	defaultValues [handlesMax]uint8
	consumer      [32]byte
	lines         uint32
	fd            int32 // C int is 32 bits even on x86_64
}

// HandleConfig is the structure to reconfigure an existing HandleRequest (not used currently, require Kernel 5.5 or later).
type HandleConfig struct {
	flags         uint32
	defaultValues [handlesMax]uint8
	padding       [4]uint32 /* padding for future use */
}

// handleData is the structure holding values for a DataLine.
//
// When getting the state of lines this contains the current state of a line.
// When setting the state of lines these should contain the desired target state.
type handleData struct {
	values [handlesMax]uint8
}

const (
	ioctlHandleGetLineValues = ((iocRead | iocWrite) << iocDirShift) | (0xB4 << iocTypeShift) | (0x08 << iocNRShift) | (unsafe.Sizeof(handleData{}) << iocSizeShift)
	ioctlHandleSetLineValues = ((iocRead | iocWrite) << iocDirShift) | (0xB4 << iocTypeShift) | (0x09 << iocNRShift) | (unsafe.Sizeof(handleData{}) << iocSizeShift)
)

// EventLineType defines the kind of event to wait on an event line.
type EventLineType uint32

// EventLine request flags.
const (
	RisingEdge  EventLineType = 1 << 0
	FallingEdge               = 1 << 1
	BothEdges                 = (1 << 0) | (1 << 1)
)

// EventLine represents a single line setup to receive GPIO events.
type EventLine struct {
	LineOffset  uint32
	HandleFlags HandleRequestFlag
	EventFlags  uint32
	Consumer    [32]byte
	Fd          int32 // C int is 32 bits even on x86_64
}

// Event types
const (
	EventRisingEdge  = 0x01
	EventFallingEdge = 0x02
)

// Event represents a occured event.
type Event struct {
	Timestamp uint64
	ID        uint32
}

const (
	ioctlGetChipInfo   = (iocRead << iocDirShift) | (0xB4 << iocTypeShift) | (0x01 << iocNRShift) | (unsafe.Sizeof(ChipInfo{}) << iocSizeShift)
	ioctlGetLineInfo   = ((iocRead | iocWrite) << iocDirShift) | (0xB4 << iocTypeShift) | (0x02 << iocNRShift) | (unsafe.Sizeof(LineInfo{}) << iocSizeShift)
	ioctlGetLineHandle = ((iocRead | iocWrite) << iocDirShift) | (0xB4 << iocTypeShift) | (0x03 << iocNRShift) | (unsafe.Sizeof(HandleRequest{}) << iocSizeShift)
	ioctlGetLineEvent  = ((iocRead | iocWrite) << iocDirShift) | (0xB4 << iocTypeShift) | (0x04 << iocNRShift) | (unsafe.Sizeof(EventLine{}) << iocSizeShift)
)
