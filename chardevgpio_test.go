// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// +build linux

package chardevgpio_test

import (
	"fmt"
	"os"
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
		t.Fatalf("Unable to open chip '%s', err='%s'", gpioDevicePath, err)
	}
	return c
}

func TestChip(t *testing.T) {
	// error cases
	_, err := gpio.NewChip("/does/not/exist")
	if err == nil {
		t.Errorf("Opening a non existing file should fail")
	} else {
		if !os.IsNotExist(err) {
			t.Errorf("Wrong err when opening a non existing file: %s", err)
		}
	}

	_, err = gpio.NewChip("/dev/zero")
	if err == nil {
		t.Errorf("Opening a invalid GPIO device should fail")
	}

	// normal case
	c := newChip(t)

	if string(c.Name()) != chipName {
		t.Errorf("Wrong value for chip name, expected='%s', got='%s'", chipName, c.Name())
	}

	if string(c.Label()) != chipLabel {
		t.Errorf("Wrong value for chip label, expected='%s', got='%s'", chipLabel, c.Label())
	}

	if c.Lines() != chipLines {
		t.Errorf("Wrong value for number of lines managed by the chip, expected=%d, got=%d", chipLines, c.Lines())
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
		if li.Offset() != i {
			t.Errorf("Wrong value for line offset, expected=%d, got=%d", i, li.Offset())
		}
		name := fmt.Sprintf("%s-%d", chipLabel, i)
		if li.Name() != name {
			t.Errorf("Wrong value for line name, expected=%s, got=%s", name, li.Name())
		}
	}
}

func TestRequestLine(t *testing.T) {
	testCases := []struct {
		offsets   []int
		direction gpio.HandleRequestFlag
		consumer  string
	}{
		{[]int{0}, gpio.HandleRequestOutput, "myapp"},
		{[]int{1}, gpio.HandleRequestInput, ""},
		{[]int{0, 1}, gpio.HandleRequestInput, "myappwithatoomanylongnamethisisnotreasonable"},
	}

	for i, tc := range testCases {
		c := newChip(t)

		l := gpio.NewHandleRequest(tc.offsets, tc.direction)
		if tc.consumer == "" {
			tc.consumer = "?"
		} else {
			l.WithConsumer(tc.consumer)
			if len(tc.consumer) > 32 {
				tc.consumer = tc.consumer[0:31] // 32 including \0
			}
		}
		if err := c.RequestLines(l); err != nil {
			t.Errorf("TestRequestLine [%02d]: unable to request line, err='%s'", i, err)
		}

		li, err := c.LineInfo(tc.offsets[0])
		if err != nil {
			t.Errorf("TestRequestLine [%02d]: unable to request line info, err='%s'", i, err)
		}
		if li.Consumer() != tc.consumer {
			t.Errorf("TestRequestLine [%02d]: wrong value for line consumer, expected=%s, got=%s", i, tc.consumer, li.Consumer())
		}

		l.Close()
		c.Close()
	}
}

func TestEventLine(t *testing.T) {
	c := newChip(t)
	defer c.Close()

	watcher, err := gpio.NewLineWatcher()
	if err != nil {
		t.Errorf("Unable to create LineWatcher, err='%s'", err)
	}

	if err := watcher.Add(c, 1, gpio.BothEdges, ""); err != nil {
		t.Errorf("Error while adding event to the LineWatcher, err='%s'", err)
	}

	if err := watcher.Close(); err != nil {
		t.Errorf("Error while closing LineWatcher, err='%s'", err)
	}
}
