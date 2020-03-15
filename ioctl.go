// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package chardevgpio

/**
 * Code in this file mimics the C code from uapi/asm-generic/ioctl.h
 * For reference, see https://elixir.bootlin.com/linux/v5.5.9/source/include/uapi/asm-generic/ioctl.h
 */

const (
	_IOC_NRBITS   = 8
	_IOC_TYPEBITS = 8

	_IOC_SIZEBITS = 14
	_IOC_DIRBITS  = 2

	_IOC_NRSHIFT   = 0
	_IOC_TYPESHIFT = _IOC_NRSHIFT + _IOC_NRBITS
	_IOC_SIZESHIFT = _IOC_TYPESHIFT + _IOC_TYPEBITS
	_IOC_DIRSHIFT  = _IOC_SIZESHIFT + _IOC_SIZEBITS

	_IOC_NONE  = 0
	_IOC_READ  = 2
	_IOC_WRITE = 1
)

func _IOR(aType, nr, size uintptr) uintptr {
	return _IOC(_IOC_READ, aType, nr, size)
}

func _IOWR(aType, nr, size uintptr) uintptr {
	return _IOC(_IOC_READ|_IOC_WRITE, aType, nr, size)
}

func _IOC(dir, aType, nr, size uintptr) uintptr {
	return (dir << _IOC_DIRSHIFT) | (aType << _IOC_TYPESHIFT) | (nr << _IOC_NRSHIFT) | (size << _IOC_SIZESHIFT)
}
