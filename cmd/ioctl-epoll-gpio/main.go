// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Inspired of https://framagit.org/cpb/ioctl-access-to-gpio/blob/master/ioctl-poll-gpio.c from Christophe Blaess.
// C version: https://gist.github.com/vinymeuh/c892df73407d0b336c879a7c87be0db7
//
// GOOS=linux GOARCH=arm GOARM=7 go build

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"

	"github.com/vinymeuh/chardevgpio"
)

func main() {
	devicePath := flag.String("device", "/dev/gpiochip0", "GPIO device path")
	lineOffset := flag.Int("line", 20, "input line number")
	flag.Parse()

	// Open the chip
	chip, err := chardevgpio.Open(*devicePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "chardevgpio.Open: %s\n", err)
		os.Exit(1)
	}
	defer chip.Close()

	// Request the EventLine
	eventLine, err := chip.RequestEventLine(*lineOffset, filepath.Base(os.Args[0]), chardevgpio.BothEdge)
	if err != nil {
		fmt.Fprintf(os.Stderr, "chip.RequestEventLine: %s\n", err)
		os.Exit(1)
	}
	defer eventLine.Close()

	// Wait for events
	var events [1]unix.EpollEvent
	for {
		nevents, err := unix.EpollWait(eventLine.EpollFd, events[:], -1)
		if err != nil {
			if err == unix.EINTR {
				continue
			}
			fmt.Fprintf(os.Stderr, "EpollWait: %s\n", err)
			os.Exit(1)
		}

		for i := 0; i < nevents; i++ {
			ev := events[i]
			if ev.Events&unix.EPOLLIN != 0 {
				evds, err := readEventsData(int(ev.Fd))
				if err != nil {
					fmt.Fprintf(os.Stderr, "readEventData: %s\n", err)
				}

				for _, evd := range evds {
					fmt.Printf("[%d.%09d]", evd.Timestamp/1000000000, evd.Timestamp%1000000000)
					if evd.Id&chardevgpio.GPIOEVENT_EVENT_RISING_EDGE == chardevgpio.GPIOEVENT_EVENT_RISING_EDGE {
						fmt.Fprintln(os.Stdout, " RISING")
					}
					if evd.Id&chardevgpio.GPIOEVENT_EVENT_FALLING_EDGE == chardevgpio.GPIOEVENT_EVENT_FALLING_EDGE {
						fmt.Fprintln(os.Stdout, " FALLING")
					}
				}
			}
		}
	}

}

// How to know that buffer size must be 16 ?
// GPIOEventData = uint64 + uint32 = 8 + 4 = 12 !?
const BUFFER_SIZE = 16

func readEventsData(fd int) ([]chardevgpio.GPIOEventData, error) {
	var evds []chardevgpio.GPIOEventData
	var evd chardevgpio.GPIOEventData
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
