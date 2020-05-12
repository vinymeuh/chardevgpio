// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// +build linux

package chardevgpio_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

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
	assert.NoErrorf(t, err, "unable to open chip")
	return c
}

func TestChip(t *testing.T) {
	// error cases
	_, err := gpio.NewChip("/does/not/exist")
	assert.Error(t, err, "opening a non existing file should fail")
	assert.Truef(t, os.IsNotExist(err), "wrong err when opening a non existing file: %s", err)

	_, err = gpio.NewChip("/dev/zero")
	assert.Error(t, err, "opening a invalid GPIO device should fail")

	// normal case
	c := newChip(t)
	assert.Equal(t, chipName, c.Name(), "wrong value for chip name")
	assert.Equal(t, chipLabel, c.Label(), "wrong value for chip label")
	assert.Equal(t, chipLines, c.Lines(), "wrong value for number of lines managed by the chip")
	assert.NoErrorf(t, c.Close(), "error while closing the chip")
	assert.Errorf(t, c.Close(), "double close a chip should return an error")
}

func TestLineInfo(t *testing.T) {
	var c gpio.Chip

	// error cases
	c = newChip(t)
	_, err := c.LineInfo(chipLines + 1)
	assert.Errorf(t, err, "requesting a LineInfo with invalide offset should fail")
	c.Close()

	// normal case
	c = newChip(t)
	for i := 0; i < c.Lines(); i++ {
		li, err := c.LineInfo(i)
		assert.NoErrorf(t, err, "unable to read LineInfo for line n°%d", i)
		assert.Equal(t, i, li.Offset(), "wrong offset for line n°%d", i)

		name := fmt.Sprintf("%s-%d", chipLabel, i)
		assert.Equal(t, name, li.Name(), "wrong name for line n°%d", i)
	}
	c.Close()
}

func TestHandleRequest(t *testing.T) {
	// prepare request for more that 64 offsets or default values
	toomany := make([]int, 128, 128)
	assert.Panics(t, func() { gpio.NewHandleRequest(toomany, gpio.HandleRequestOutput) })
	assert.Panics(t, func() { gpio.NewHandleRequest([]int{0}, gpio.HandleRequestOutput).WithDefaults(toomany) })
}

func TestRequestLineBusy(t *testing.T) {
	c := newChip(t)

	li := gpio.NewHandleRequest([]int{0}, gpio.HandleRequestOutput)
	c.RequestLines(li)
	assert.Errorf(t, c.RequestLines(li), "should have return a 'device or resource busy' error")
	li.Close()

	c.Close()
}

func TestRequestLine(t *testing.T) {
	testCases := []struct {
		offsets   []int
		direction gpio.HandleRequestFlag
		consumer  string
		defaults  []int
	}{
		{[]int{0}, gpio.HandleRequestOutput, "myapp", []int{}},
		{[]int{0}, gpio.HandleRequestInput, "myappwithatoomanylongnamethisisnotreasonable", []int{}},
		{[]int{0, 1}, gpio.HandleRequestInput, "", []int{}},
		{[]int{0, 1}, gpio.HandleRequestOutput, "", []int{1, 1}},
		{[]int{0}, gpio.HandleRequestOutput | gpio.HandleRequestActiveLow, "", []int{}},
		{[]int{0}, gpio.HandleRequestOutput | gpio.HandleRequestOpenDrain, "", []int{}},
		{[]int{0}, gpio.HandleRequestOutput | gpio.HandleRequestOpenSource, "", []int{}},
	}

	for i, tc := range testCases {
		// setup request
		l := gpio.NewHandleRequest(tc.offsets, tc.direction)
		if tc.consumer == "" {
			tc.consumer = "?"
		} else {
			l.WithConsumer(tc.consumer)
			if len(tc.consumer) > 32 {
				tc.consumer = tc.consumer[0:31] // 32 including \0
			}
		}
		if len(tc.defaults) > 0 {
			l.WithDefaults(tc.defaults)
		}

		// tests
		c := newChip(t)
		err := c.RequestLines(l)
		assert.NoErrorf(t, err, "Test n°%02d, unable to request line, err='%s'", i, err)

		li, err := c.LineInfo(tc.offsets[0])
		assert.NoErrorf(t, err, "Test n°%02d, unable to request line info, err='%s'", i, err)
		assert.Equal(t, tc.consumer, li.Consumer(), "Test n°%02d, wrong value for line consumer", i)

		switch tc.direction {
		case tc.direction & gpio.HandleRequestInput:
			assert.Truef(t, li.IsInput(), "Test n°%02d, should be an input line", i)
		case tc.direction & gpio.HandleRequestOutput:
			assert.Truef(t, li.IsOutput(), "Test n°%02d, should be an output line", i)
		}

		if tc.direction&gpio.HandleRequestActiveLow == gpio.HandleRequestActiveLow {
			assert.Truef(t, li.IsActiveLow(), "Test n°%02d, should be active low", i)
		} else {
			assert.Truef(t, li.IsActiveHigh(), "Test n°%02d, should be active high", i)
		}

		if tc.direction&gpio.HandleRequestOpenDrain == gpio.HandleRequestOpenDrain {
			assert.Truef(t, li.IsOpenDrain(), "Test n°%02d, should be open drain", i)
		} else {
			assert.Falsef(t, li.IsOpenDrain(), "Test n°%02d, should not be open drain", i)
		}

		if tc.direction&gpio.HandleRequestOpenSource == gpio.HandleRequestOpenSource {
			assert.Truef(t, li.IsOpenSource(), "Test n°%02d, should be open source", i)
		} else {
			assert.Falsef(t, li.IsOpenSource(), "Test n°%02d, should not be open source", i)
		}

		// teardown
		assert.NoErrorf(t, l.Close(), "Test n°%02d, error while closing the line", i)
		c.Close()
	}
}

func TestRequestLineErrOperationNotPermitted(t *testing.T) {
	c := newChip(t)

	l := gpio.NewHandleRequest([]int{0}, gpio.HandleRequestOutput)
	c.RequestLines(l)
	_, _, err := l.Read()
	assert.Error(t, err)
	assert.Equal(t, gpio.ErrOperationNotPermitted, err)
	l.Close()

	l = gpio.NewHandleRequest([]int{1}, gpio.HandleRequestInput)
	c.RequestLines(l)
	err = l.Write(1)
	assert.Error(t, err)
	assert.Equal(t, gpio.ErrOperationNotPermitted, err)
	l.Close()	

	c.Close()
}

func TestEventLine(t *testing.T) {
	c := newChip(t)

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

	c.Close()
}
