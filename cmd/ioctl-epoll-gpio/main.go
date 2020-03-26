// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Inspired of https://framagit.org/cpb/ioctl-access-to-gpio/blob/master/ioctl-poll-gpio.c from Christophe Blaess.
// C version: https://gist.github.com/vinymeuh/c892df73407d0b336c879a7c87be0db7
//
// GOOS=linux GOARCH=arm GOARM=7 go build

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vinymeuh/chardevgpio"
)

func printEventData(evd chardevgpio.GPIOEventData) {
	fmt.Printf("[%d.%09d]", evd.Timestamp/1000000000, evd.Timestamp%1000000000)
	if evd.Id&chardevgpio.GPIOEVENT_EVENT_RISING_EDGE == chardevgpio.GPIOEVENT_EVENT_RISING_EDGE {
		fmt.Fprintln(os.Stdout, " RISING")
	}
	if evd.Id&chardevgpio.GPIOEVENT_EVENT_FALLING_EDGE == chardevgpio.GPIOEVENT_EVENT_FALLING_EDGE {
		fmt.Fprintln(os.Stdout, " FALLING")
	}
}

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

	// Create the EventLineWatcher
	watcher, err := chardevgpio.NewEventLineWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "chardevgpio.NewEventLineWatcher: %s\n", err)
		os.Exit(1)
	}
	defer watcher.Close()

	if err := watcher.AddEvent(chip, *lineOffset, filepath.Base(os.Args[0]), chardevgpio.BothEdge); err != nil {
		fmt.Fprintf(os.Stderr, "watcher.AddEvent: %s\n", err)
		os.Exit(1)
	}

	err = watcher.WaitForEver(printEventData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "watcher.WaitForEver: %s\n", err)
		os.Exit(1)
	}
}
