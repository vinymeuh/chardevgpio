// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// +build linux

package chardevgpio

import "unsafe"

/**
 * Code in this file mimics directly the Linux kernel code.
 *
 * For reference, see:
 *  - https://elixir.bootlin.com/linux/v5.5.9/source/include/uapi/asm-generic/ioctl.h
 *  - https://elixir.bootlin.com/linux/v5.5.9/source/include/uapi/linux/gpio.h
 */

const (
	iocNRBITS   = 8
	iocTYPEBITS = 8

	iocSIZEBITS = 14
	iocDIRBITS  = 2

	iocNRSHIFT   = 0
	iocTYPESHIFT = iocNRSHIFT + iocNRBITS
	iocSIZESHIFT = iocTYPESHIFT + iocTYPEBITS
	iocDIRSHIFT  = iocSIZESHIFT + iocSIZEBITS

	iocNONE  = 0
	iocREAD  = 2
	iocWRITE = 1
)

// ChipInfo contains informations about a GPIO chip.
type ChipInfo struct {
	Name  [32]byte
	Label [32]byte
	Lines uint32
}

// Informational flags
const (
	gpioLineFlagKernel       = 1 << 0
	gpioLineFlagIsOut        = 1 << 1
	gpioLineFlagActiveLow    = 1 << 2
	gpioLineFlagOpenDrain    = 1 << 3
	gpioLineFlagOpenSource   = 1 << 4
	gpioLineFlagBiasPullUp   = 1 << 5
	gpioLineFlagBiasPullDown = 1 << 6
	gpioLineFlagDisable      = 1 << 7
)

// LineInfo contains informations about a GPIO line.
type LineInfo struct {
	Offset   uint32
	Flags    uint32
	Name     [32]byte
	Consumer [32]byte
}

// gpioHandlesMax limits maximum number of handles that can be requested in a GPIOHandleRequest
const gpioHandlesMax = 64

// DataLine request flags.
const (
	gpioHandleRequestInput        = 1 << 0
	gpioHandleRequestOutput       = 1 << 1
	gpioHandleRequestActiveLow    = 1 << 2
	gpioHandleRequestOpenDrain    = 1 << 3
	gpioHandleRequestOpenSource   = 1 << 4
	gpioHandleRequestBiasPullUp   = 1 << 5
	gpioHandleRequestBiasPullDown = 1 << 6
	gpioHandleRequestBiasDisable  = 1 << 7
)

// DataLine represents a single line to be used to send or receive data.
type DataLine struct {
	LineOffsets   [gpioHandlesMax]uint32
	Flags         uint32
	DefaultValues [gpioHandlesMax]uint8
	Consumer      [32]byte
	Lines         uint32
	Fd            int
}

// DataLineConfig is the structure to configure a DataLine.
type DataLineConfig struct {
	Flags         uint32
	DefaultValues [gpioHandlesMax]uint8
	Padding       [4]uint32 /* padding for future use */
}

// Data is the structure holding values for a DataLine.
//
// When getting the state of lines this contains the current state of a line.
// When setting the state of lines these should contain the desired target state.
type Data struct {
	Values [gpioHandlesMax]uint8
}

const (
	gpioHandleGetLineValuesIOCTL = ((iocREAD | iocWRITE) << iocDIRSHIFT) | (0xB4 << iocTYPESHIFT) | (0x08 << iocNRSHIFT) | (unsafe.Sizeof(Data{}) << iocSIZESHIFT)
	gpioHandleSetLineValuesIOCTL = ((iocREAD | iocWRITE) << iocDIRSHIFT) | (0xB4 << iocTYPESHIFT) | (0x09 << iocNRSHIFT) | (unsafe.Sizeof(Data{}) << iocSIZESHIFT)
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
	HandleFlags uint32
	EventFlags  uint32
	Consumer    [32]byte
	Fd          int
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
	gpioGetChipInfoIOCTL   = (iocREAD << iocDIRSHIFT) | (0xB4 << iocTYPESHIFT) | (0x01 << iocNRSHIFT) | (unsafe.Sizeof(ChipInfo{}) << iocSIZESHIFT)
	gpioGetLineInfoIOCTL   = ((iocREAD | iocWRITE) << iocDIRSHIFT) | (0xB4 << iocTYPESHIFT) | (0x02 << iocNRSHIFT) | (unsafe.Sizeof(LineInfo{}) << iocSIZESHIFT)
	gpioGetLineHandleIOCTL = ((iocREAD | iocWRITE) << iocDIRSHIFT) | (0xB4 << iocTYPESHIFT) | (0x03 << iocNRSHIFT) | (unsafe.Sizeof(DataLine{}) << iocSIZESHIFT)
	gpioGetLineEventIOCTL  = ((iocREAD | iocWRITE) << iocDIRSHIFT) | (0xB4 << iocTYPESHIFT) | (0x04 << iocNRSHIFT) | (unsafe.Sizeof(EventLine{}) << iocSIZESHIFT)
)
