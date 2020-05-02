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

func main() {
	path := flag.String("device", "/dev/gpiochip0", "GPIO device path")
	offset := flag.Int("line", 22, "line number")
	flag.Parse()

	chip, err := gpio.NewChip(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer chip.Close()

	line := gpio.NewHandleRequest([]int{*offset}, gpio.HandleRequestInput).WithConsumer(filepath.Base(os.Args[0]))
	if err := chip.RequestLines(line); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer line.Close()

	value, _, err := line.Read()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(value)
}
