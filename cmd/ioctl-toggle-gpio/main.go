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
	"time"

	"github.com/vinymeuh/chardevgpio"
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

	outputLine, err := chip.RequestOutputLine(*lineOffset, filepath.Base(os.Args[0]))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer outputLine.Close()

	value := 0
	for {
		value = 1 - value
		err := outputLine.SetValue(value)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		time.Sleep(500 * time.Millisecond)
	}

}
