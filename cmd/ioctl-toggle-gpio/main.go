// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Inspired of https://framagit.org/cpb/ioctl-access-to-gpio/blob/master/ioctl-toggle-gpio.c from Christophe Blaess.
//
// GOOS=linux GOARCH=arm GOARM=7 go build

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"

	"github.com/vinymeuh/chardevgpio"
	"golang.org/x/sys/unix"
)

func main() {
	devicePath := flag.String("device", "/dev/gpiochip0", "GPIO device path")
	lineOffset := flag.Int("line", 22, "line number")
	flag.Parse()

	chip, err := chardevgpio.Open(*devicePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer chip.Close()

	outputRequest := chardevgpio.GPIOHandleRequest{}
	outputRequest.LineOffsets[0] = uint32(*lineOffset)
	outputRequest.Flags = chardevgpio.GPIOHANDLE_REQUEST_OUTPUT
	outputRequest.Lines = uint32(1)

	var consumer [32]byte
	for i, c := range []byte(filepath.Base(os.Args[0])) {
		if i == 32 {
			break
		}
		consumer[i] = c
	}
	outputRequest.Consumer = consumer

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, chip.Fd, chardevgpio.GPIO_GET_LINEHANDLE_IOCTL, uintptr(unsafe.Pointer(&outputRequest)))
	if errno != 0 {
		fmt.Fprintln(os.Stderr, errno)
		os.Exit(1)
	}

	outputValue := chardevgpio.GPIOHandleData{}
	outputValue.Values[0] = 0
	for {
		outputValue.Values[0] = 1 - outputValue.Values[0]
		_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(outputRequest.Fd), chardevgpio.GPIOHANDLE_SET_LINE_VALUES_IOCTL, uintptr(unsafe.Pointer(&outputValue)))
		if errno != 0 {
			fmt.Fprintln(os.Stderr, errno)
			syscall.Close(int(outputRequest.Fd))
			os.Exit(1)
		}
		time.Sleep(500 * time.Millisecond)
	}

}
