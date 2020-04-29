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

	for i := 0; i < c.Lines(); i++ {
		li, err := c.LineInfo(i)
		if err != nil {
			t.Errorf("Unable to read LineInfo for line %d, err='%s'", i, err)
		}
		name := fmt.Sprintf("%s-%d", chipLabel, i)
		if string(li.Name[:len(name)]) != name {
			t.Errorf("Wrong value for LineInfo.Name, expected=%s, got=%s", name, li.Name)
		}
	}

	if err := c.Close(); err != nil {
		t.Errorf("Error while closing the chip, err='%s'", err)
	}
}
