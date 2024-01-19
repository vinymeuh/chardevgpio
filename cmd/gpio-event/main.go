// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	gpio "github.com/vinymeuh/chardevgpio"
)

func printEventData(evd gpio.Event) {
	fmt.Printf("[%d.%09d]", evd.Timestamp/1000000000, evd.Timestamp%1000000000)
	if evd.IsRising() {
		fmt.Fprintln(os.Stdout, " RISING")
	}
	if evd.IsFalling() {
		fmt.Fprintln(os.Stdout, " FALLING")
	}
}

func main() {
	devicePath := flag.String("device", "/dev/gpiochip0", "GPIO device path")
	lineOffset := flag.Int("line", 20, "input line number")
	flag.Parse()

	// Open the chip
	chip, err := gpio.NewChip(*devicePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "chardevgpio.Open: %s\n", err)
		os.Exit(1)
	}
	defer chip.Close()

	// Create the LineWatcher
	watcher, err := gpio.NewLineWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "chardevgpio.NewEventLineWatcher: %s\n", err)
		os.Exit(1)
	}
	defer watcher.Close()

	if err := watcher.Add(chip, *lineOffset, gpio.BothEdges, filepath.Base(os.Args[0])); err != nil {
		fmt.Fprintf(os.Stderr, "watcher.AddEvent: %s\n", err)
		os.Exit(1)
	}

	err = watcher.WaitForEver(printEventData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "watcher.WaitForEver: %s\n", err)
		os.Exit(1)
	}
}
