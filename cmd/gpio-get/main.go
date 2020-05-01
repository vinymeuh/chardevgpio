// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	gpio "github.com/vinymeuh/chardevgpio"
)

func main() {
	path := flag.String("device", "/dev/gpiochip0", "GPIO device path")
	offset := flag.Int("line", 22, "line number")
	value := flag.Int("value", 1, "value to write (0/1)")
	seconds := flag.Int("time", 60, "write hold time (seconds)")
	flag.Parse()

	chip, err := gpio.NewChip(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer chip.Close()

	line := gpio.NewHandleRequest([]int{*offset}, gpio.HandleRequestInput)
	if err := chip.RequestLines(line); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer line.Close()

	if err := line.Write0(*value); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	time.Sleep(time.Duration(*seconds) * time.Second)
}
