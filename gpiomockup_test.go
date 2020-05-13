// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// +build linux

package chardevgpio_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var mockChip MockChip

func init() {
	var err error
	mockChip, err = NewMockChip()
	if err != nil {
		fmt.Printf("Unable to initialize MockChip: %s\n", err)
		os.Exit(1)
	}
}

type MockChip struct {
	Path      string
	Name      string
	Label     string
	Lines     int
	debugPath string
}

func NewMockChip() (MockChip, error) {
	c := MockChip{
		Path:      "/dev/gpiochip0",
		Name:      "gpiochip0",
		Label:     "gpio-mockup-A",
		Lines:     10,
		debugPath: "/sys/kernel/debug/gpio-mockup/gpiochip0",
	}

	_, err := os.Stat(c.debugPath)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (m MockChip) linePath(i int) string {
	return fmt.Sprintf("%s/%d", m.debugPath, i)
}

func (m MockChip) Read() ([]int, error) {
	out := make([]int, m.Lines)
	for i := 0; i < m.Lines; i++ {
		cmdStr := fmt.Sprintf("cat %s", m.linePath(i))
		cmd := exec.Command("sh", "-c", cmdStr)
		var cmdOut bytes.Buffer
		cmd.Stdout = &cmdOut
		err := cmd.Run()
		if err != nil {
			return out, err
		}

		out[i], err = strconv.Atoi(strings.TrimSpace(cmdOut.String()))
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

func (m MockChip) Write(data []int) error {
	for i := range data {
		cmdStr := fmt.Sprintf("echo \"%d\" > %s", data[i], m.linePath(i))
		cmd := exec.Command("sh", "-c", cmdStr)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
