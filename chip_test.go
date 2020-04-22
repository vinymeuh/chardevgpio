// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// +build linux

// Package chardevgpio is a low-level library to the Linux GPIO Character device API.
package chardevgpio

import (
	"fmt"
	"testing"
)

const (
	gpioDevicePath = "/dev/gpiochip0"
	chipName       = "gpiochip0"
	chipLabel      = "gpio-mockup-A"
	chipLines      = 10
)

func TestChip(t *testing.T) {
	chip, err := Open(gpioDevicePath)
	if err != nil {
		t.Fatalf("Unable to open gpio device '%s', err='%s'", gpioDevicePath, err)
	}

	// we prefere to fail if we are not sure we are using the expected gpio-mockup device
	if string(chip.Label[:len(chipLabel)]) != chipLabel {
		t.Fatalf("Wrong value for Chip.Label, expected='%s', got='%s'", chipLabel, chip.Label)
	}

	if string(chip.Name[:len(chipName)]) != chipName {
		t.Errorf("Wrong value for Chip.Name, expected='%s', got='%s'", chipName, chip.Name)
	}

	if chip.Lines != chipLines {
		t.Errorf("Wrong value for Chip.Lines, expected=%d, got=%d", chipLines, chip.Lines)
	}

	for i := 0; i < int(chip.Lines); i++ {
		li, err := chip.LineInfo(i)
		if err != nil {
			t.Errorf("Unable to read LineInfo for line %d, err='%s'", i, err)
		}
		name := fmt.Sprintf("%s-%d", chipLabel, i)
		if string(li.Name[:len(name)]) != name {
			t.Errorf("Wrong value for LineInfo.Name, expected=%s, got=%s", name, li.Name)
		}
	}

	if err := chip.Close(); err != nil {
		t.Errorf("Error while closing the chip, err='%s'", err)
	}
}
