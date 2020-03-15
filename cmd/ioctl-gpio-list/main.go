// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Inspired of https://framagit.org/cpb/ioctl-access-to-gpio/blob/master/ioctl-gpio-list.c from Christophe Blaess.
//
// GOOS=linux GOARCH=arm GOARM=7 go build

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vinymeuh/chardevgpio"
)

func printChipInfo(path string) {
	chip, err := chardevgpio.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer chip.Close()
	fmt.Printf("file = %s, name = %s, label = %s, lines = %d\n", path, chip.Name, chip.Label, chip.Lines)

	for i := 0; i < int(chip.Lines); i++ {
		l, err := chip.GetLine(i)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		printLineInfo(l)
	}
}

func printLineInfo(l chardevgpio.Line) {

	fmt.Printf("    line %2d: name = \"%s\", consumer = \"%s\", flags = ", l.Offset, l.Name, l.Consumer)
	if l.IsOut() {
		fmt.Print("OUT")
	} else {
		fmt.Print("IN ")
	}
	if l.IsActiveLow() {
		fmt.Print(" ACTIVE_LOW ")
	} else {
		fmt.Print(" ACTIVE_HIGH")
	}
	if l.IsOpenDrain() {
		fmt.Print(" OPEN_DRAIN")
	}
	if l.IsOpenSource() {
		fmt.Print(" OPEN_SOURCE")
	}
	if l.IsKernel() {
		fmt.Print(" KERNEL")
	}

	fmt.Println()
}

func main() {
	chips, _ := filepath.Glob("/dev/gpiochip*")
	for _, chip := range chips {
		printChipInfo(chip)
	}
}
