// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// +build linux

// Package chardevgpio is a low-level library to the Linux GPIO Character device API.
package chardevgpio_test

import (
	"testing"

	"github.com/vinymeuh/chardevgpio"
)

func TestOutputDataLine(t *testing.T) {
	chip, err := chardevgpio.NewChip(gpioDevicePath)
	if err != nil {
		t.Fatalf("Unable to open gpio device '%s', err='%s'", gpioDevicePath, err)
	}

	_, err = chip.RequestOutputLine(0, 0, "TestOutputDataLine")
	if err != nil {
		t.Fatalf("Unable to get line 0 for output, err='%s'", err)
	}

	chip.Close()
}

func TestInputDataLine(t *testing.T) {
	chip, err := chardevgpio.NewChip(gpioDevicePath)
	if err != nil {
		t.Fatalf("Unable to open gpio device '%s', err='%s'", gpioDevicePath, err)
	}

	_, err = chip.RequestInputLine(1, "TestInputDataLine")
	if err != nil {
		t.Fatalf("Unable to get line 1 for input, err='%s'", err)
	}

	chip.Close()
}

func TestEventLine(t *testing.T) {
	chip, err := chardevgpio.NewChip(gpioDevicePath)
	if err != nil {
		t.Fatalf("Unable to open gpio device '%s', err='%s'", gpioDevicePath, err)
	}

	watcher, err := chardevgpio.NewEventLineWatcher()
	if err != nil {
		t.Fatalf("Unable to create EventLineWatcher: %s\n", err)
	}
	defer watcher.Close()

	if err := watcher.AddEvent(chip, 2, "TestEventLine", chardevgpio.BothEdges); err != nil {
		t.Fatalf("Unable to add line to watcher: %s\n", err)
	}
}
