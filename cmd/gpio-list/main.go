// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Inspired of https://framagit.org/cpb/ioctl-access-to-gpio/blob/master/ioctl-gpio-list.c from Christophe Blaess.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	gpio "github.com/vinymeuh/chardevgpio"
)

func printChipInfo(path string) {
	chip, err := gpio.NewChip(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer chip.Close()
	fmt.Printf("file = %s, name = %s, label = %s, lines = %d\n", path, chip.Name(), chip.Label(), chip.Lines())

	for i := 0; i < chip.Lines(); i++ {
		li, err := chip.LineInfo(i)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		printLineInfo(li)
	}
}

func printLineInfo(li gpio.LineInfo) {

	fmt.Printf("    line %2d: name = \"%s\", consumer = \"%s\", flags = ", li.Offset(), li.Name(), li.Consumer())
	if li.IsOutput() {
		fmt.Print("OUT")
	} else {
		fmt.Print("IN ")
	}
	if li.IsActiveLow() {
		fmt.Print(" ACTIVE_LOW ")
	} else {
		fmt.Print(" ACTIVE_HIGH")
	}
	if li.IsOpenDrain() {
		fmt.Print(" OPEN_DRAIN")
	}
	if li.IsOpenSource() {
		fmt.Print(" OPEN_SOURCE")
	}
	if li.IsKernel() {
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
