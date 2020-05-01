// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// +build linux

package chardevgpio_test

import (
	"fmt"
	"testing"

	gpio "github.com/vinymeuh/chardevgpio"
)

const (
	gpioDevicePath = "/dev/gpiochip0"
	chipName       = "gpiochip0"
	chipLabel      = "gpio-mockup-A"
	chipLines      = 10
)

func newChip(t *testing.T) gpio.Chip {
	c, err := gpio.NewChip(gpioDevicePath)
	if err != nil {
		t.Fatalf("Unable to open gpio device '%s', err='%s'", gpioDevicePath, err)
	}
	return c
}

func TestChip(t *testing.T) {
	c, err := gpio.NewChip(gpioDevicePath)
	if err != nil {
		t.Fatalf("Unable to open gpio device '%s', err='%s'", gpioDevicePath, err)
	}

	if string(c.Name()) != chipName {
		t.Errorf("Wrong value for Chip.Name, expected='%s', got='%s'", chipName, c.Name())
	}

	if string(c.Label()) != chipLabel {
		t.Errorf("Wrong value for Chip.Label, expected='%s', got='%s'", chipLabel, c.Label())
	}

	if c.Lines() != chipLines {
		t.Errorf("Wrong value for Chip.Lines, expected=%d, got=%d", chipLines, c.Lines())
	}

	if err := c.Close(); err != nil {
		t.Errorf("Error while closing the chip, err='%s'", err)
	}
}

func TestLineInfo(t *testing.T) {
	c := newChip(t)
	defer c.Close()

	for i := 0; i < c.Lines(); i++ {
		li, err := c.LineInfo(i)
		if err != nil {
			t.Errorf("Unable to read LineInfo for line %d, err='%s'", i, err)
		}
		name := fmt.Sprintf("%s-%d", chipLabel, i)
		if li.Name() != name {
			t.Errorf("Wrong value for LineInfo.Name, expected=%s, got=%s", name, li.Name())
		}
	}
}

func TestRequestLine(t *testing.T) {
	c := newChip(t)
	defer c.Close()

	lines := gpio.NewHandleRequest([]int{0}, gpio.HandleRequestOutput)
	if err := c.RequestLines(lines); err != nil {
		t.Errorf("Unable to get line 0 for output, err='%s'", err)
	}
}
