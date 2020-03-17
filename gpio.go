// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package chardevgpio

import (
	"unsafe"
)

/**
 * Code in this file mimics the C code from uapi/linux/gpio.h
 * For reference, see https://elixir.bootlin.com/linux/v5.5.9/source/include/uapi/linux/gpio.h
 */

type GPIOChipInfo struct {
	Name  [32]byte
	Label [32]byte
	Lines uint32
}

/* Informational flags */
const (
	GPIOLINE_FLAG_KERNEL         = 1 << 0
	GPIOLINE_FLAG_IS_OUT         = 1 << 1
	GPIOLINE_FLAG_ACTIVE_LOW     = 1 << 2
	GPIOLINE_FLAG_OPEN_DRAIN     = 1 << 3
	GPIOLINE_FLAG_OPEN_SOURCE    = 1 << 4
	GPIOLINE_FLAG_BIAS_PULL_UP   = 1 << 5
	GPIOLINE_FLAG_BIAS_PULL_DOWN = 1 << 6
	GPIOLINE_FLAG_BIAS_DISABLE   = 1 << 7
)

type GPIOLineInfo struct {
	Offset   uint32
	Flags    uint32
	Name     [32]byte
	Consumer [32]byte
}

/* Maximum number of requested handles */
const GPIOHANDLES_MAX = 64

/* Linerequest flags */
const (
	GPIOHANDLE_REQUEST_INPUT          = 1 << 0
	GPIOHANDLE_REQUEST_OUTPUT         = 1 << 1
	GPIOHANDLE_REQUEST_ACTIVE_LOW     = 1 << 2
	GPIOHANDLE_REQUEST_OPEN_DRAIN     = 1 << 3
	GPIOHANDLE_REQUEST_OPEN_SOURCE    = 1 << 4
	GPIOHANDLE_REQUEST_BIAS_PULL_UP   = 1 << 5
	GPIOHANDLE_REQUEST_BIAS_PULL_DOWN = 1 << 6
	GPIOHANDLE_REQUEST_BIAS_DISABLE   = 1 << 7
)

type GPIOHandleRequest struct {
	LineOffsets   [GPIOHANDLES_MAX]uint32
	Flags         uint32
	DefaultValues [GPIOHANDLES_MAX]uint8
	Consumer      [32]byte
	Lines         uint32
	Fd            int
}

type GPIOHandleConfig struct {
	Flags         uint32
	DefaultValues [GPIOHANDLES_MAX]uint8
	Padding       [4]uint32 /* padding for future use */
}

type GPIOHandleData struct {
	Values [GPIOHANDLES_MAX]uint8
}

/**
 * GPIOHANDLE_GET_LINE_VALUES_IOCTL _IOWR(0xB4, 0x08, unsafe.Sizeof(GPIOHandleData))
 * GPIOHANDLE_SET_LINE_VALUES_IOCTL _IOWR(0xB4, 0x09, unsafe.Sizeof(GPIOHandleData))
 */
const (
	GPIOHANDLE_GET_LINE_VALUES_IOCTL = ((_IOC_READ | _IOC_WRITE) << _IOC_DIRSHIFT) | (0xB4 << _IOC_TYPESHIFT) | (0x08 << _IOC_NRSHIFT) | (unsafe.Sizeof(GPIOHandleData{}) << _IOC_SIZESHIFT)
	GPIOHANDLE_SET_LINE_VALUES_IOCTL = ((_IOC_READ | _IOC_WRITE) << _IOC_DIRSHIFT) | (0xB4 << _IOC_TYPESHIFT) | (0x09 << _IOC_NRSHIFT) | (unsafe.Sizeof(GPIOHandleData{}) << _IOC_SIZESHIFT)
)

/**
 * GPIO_GET_CHIPINFO_IOCTL = _IOR(0xB4, 0x01, unsafe.Sizeof(GPIOChipInfo{}))
 * GPIO_GET_LINEINFO_IOCTL = _IOWR(0xB4, 0x02, unsafe.Sizeof(GPIOLineInfo{}))
 * GPIO_GET_LINEHANDLE_IOCTL = _IOWR(0xB4, 0x03, unsafe.Sizeof(GPIOHandleRequest))
 */
const (
	GPIO_GET_CHIPINFO_IOCTL   = (_IOC_READ << _IOC_DIRSHIFT) | (0xB4 << _IOC_TYPESHIFT) | (0x01 << _IOC_NRSHIFT) | (unsafe.Sizeof(GPIOChipInfo{}) << _IOC_SIZESHIFT)
	GPIO_GET_LINEINFO_IOCTL   = ((_IOC_READ | _IOC_WRITE) << _IOC_DIRSHIFT) | (0xB4 << _IOC_TYPESHIFT) | (0x02 << _IOC_NRSHIFT) | (unsafe.Sizeof(GPIOLineInfo{}) << _IOC_SIZESHIFT)
	GPIO_GET_LINEHANDLE_IOCTL = ((_IOC_READ | _IOC_WRITE) << _IOC_DIRSHIFT) | (0xB4 << _IOC_TYPESHIFT) | (0x03 << _IOC_NRSHIFT) | (unsafe.Sizeof(GPIOHandleRequest{}) << _IOC_SIZESHIFT)
)
