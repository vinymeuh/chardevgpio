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

/**
 * GPIO_GET_CHIPINFO_IOCTL = _IOR(0xB4, 0x01, unsafe.Sizeof(GPIOChipInfo{}))
 * GPIO_GET_LINEINFO_IOCTL = _IOWR(0xB4, 0x02, unsafe.Sizeof(GPIOLineInfo{}))
 */
const (
	GPIO_GET_CHIPINFO_IOCTL = (_IOC_READ << _IOC_DIRSHIFT) | (0xB4 << _IOC_TYPESHIFT) | (0x01 << _IOC_NRSHIFT) | (unsafe.Sizeof(GPIOChipInfo{}) << _IOC_SIZESHIFT)
	GPIO_GET_LINEINFO_IOCTL = ((_IOC_READ | _IOC_WRITE) << _IOC_DIRSHIFT) | (0xB4 << _IOC_TYPESHIFT) | (0x02 << _IOC_NRSHIFT) | (unsafe.Sizeof(GPIOLineInfo{}) << _IOC_SIZESHIFT)
)
