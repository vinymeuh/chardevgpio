// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// +build linux

package chardevgpio_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	gpio "github.com/vinymeuh/chardevgpio"
)

func newChip(t *testing.T) gpio.Chip {
	c, err := gpio.NewChip(mockChip.Path)
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
	assert.Equal(t, mockChip.Name, c.Name(), "wrong value for chip name")
	assert.Equal(t, mockChip.Label, c.Label(), "wrong value for chip label")
	assert.Equal(t, mockChip.Lines, c.Lines(), "wrong value for number of lines managed by the chip")
	assert.NoErrorf(t, c.Close(), "error while closing the chip")
	assert.Errorf(t, c.Close(), "double close a chip should return an error")
}

func TestLineInfo(t *testing.T) {
	var c gpio.Chip

	// error cases
	c = newChip(t)
	_, err := c.LineInfo(mockChip.Lines + 1)
	assert.Errorf(t, err, "requesting a LineInfo with invalide offset should fail")
	c.Close()

	// normal case
	c = newChip(t)
	for i := 0; i < c.Lines(); i++ {
		li, err := c.LineInfo(i)
		assert.NoErrorf(t, err, "unable to read LineInfo for line n°%d", i)
		assert.Equal(t, i, li.Offset(), "wrong offset for line n°%d", i)

		name := fmt.Sprintf("%s-%d", mockChip.Label, i)
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
		assert.NoErrorf(t, err, "test n°%02d, unable to request line, err='%s'", i, err)

		li, err := c.LineInfo(tc.offsets[0])
		assert.NoErrorf(t, err, "test n°%02d, unable to request line info, err='%s'", i, err)
		assert.Equal(t, tc.consumer, li.Consumer(), "Test n°%02d, wrong value for line consumer", i)

		switch tc.direction {
		case tc.direction & gpio.HandleRequestInput:
			assert.Truef(t, li.IsInput(), "test n°%02d, should be an input line", i)
		case tc.direction & gpio.HandleRequestOutput:
			assert.Truef(t, li.IsOutput(), "test n°%02d, should be an output line", i)
		}

		if tc.direction&gpio.HandleRequestActiveLow == gpio.HandleRequestActiveLow {
			assert.Truef(t, li.IsActiveLow(), "Test n°%02d, should be active low", i)
		} else {
			assert.Truef(t, li.IsActiveHigh(), "test n°%02d, should be active high", i)
		}

		if tc.direction&gpio.HandleRequestOpenDrain == gpio.HandleRequestOpenDrain {
			assert.Truef(t, li.IsOpenDrain(), "test n°%02d, should be open drain", i)
		} else {
			assert.Falsef(t, li.IsOpenDrain(), "test n°%02d, should not be open drain", i)
		}

		if tc.direction&gpio.HandleRequestOpenSource == gpio.HandleRequestOpenSource {
			assert.Truef(t, li.IsOpenSource(), "test n°%02d, should be open source", i)
		} else {
			assert.Falsef(t, li.IsOpenSource(), "test n°%02d, should not be open source", i)
		}

		// teardown
		assert.NoErrorf(t, l.Close(), "test n°%02d, error while closing the line", i)
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

func TestRequestLineInputRead(t *testing.T) {
	testCases := []struct {
		data []int
	}{
		{[]int{1, 1, 1}},
		{[]int{0, 0, 0}},
		{[]int{1, 0, 1}},
	}

	c := newChip(t)
	l := gpio.NewHandleRequest([]int{0, 1, 2}, gpio.HandleRequestInput)
	c.RequestLines(l)

	// Read
	for n, tc := range testCases {
		err := mockChip.Write(tc.data)
		assert.NoError(t, err, "unable to write using mockChip")

		_, read, err := l.Read()
		assert.NoError(t, err, "unable to read from input line")
		for i := range tc.data {
			assert.Equal(t, tc.data[i], read[i], "test n°%02d, value n°%d does not match", n, i)
		}
	}

	l.Close()
	c.Close()
}

func TestRequestLineOutputWrite(t *testing.T) {
	testCases := []struct {
		data []int
	}{
		{[]int{1, 1, 1}},
		{[]int{0, 0, 0}},
		{[]int{1, 0, 1}},
	}

	c := newChip(t)
	defaults := []int{1, 1, 1}
	l := gpio.NewHandleRequest([]int{0, 1, 2}, gpio.HandleRequestOutput).WithDefaults(defaults)
	c.RequestLines(l)

	// WithDefaults
	mockread, err := mockChip.Read()
	assert.NoError(t, err, "unable to read using mockChip")
	for i := range defaults {
		assert.Equal(t, defaults[i], mockread[i], "default value n°%d does not match", i)
	}

	// Write
	for n, tc := range testCases {
		err := l.Write(tc.data[0], tc.data[1:]...)
		assert.NoError(t, err, "unable to write to output line")

		mockread, err := mockChip.Read()
		assert.NoError(t, err, "unable to read using mockChip")
		for i := range tc.data {
			assert.Equal(t, tc.data[i], mockread[i], "test n°%02d, value n°%d does not match", n, i)
		}
	}

	l.Close()
	c.Close()
}

func TestEventLine(t *testing.T) {
	done := make(chan struct{}, 1)
	line := 0

	go testEventLineWait(t, line, done)
	mockChip.Write([]int{0})
	mockChip.Write([]int{1})

	select {
	case <-done:
		break
	case <-time.After(2 * time.Second):
		assert.Fail(t, "testEventLineWait did not finished before timeout")
	}
}

func testEventLineWait(t *testing.T, line int, done chan struct{}) {
	c := newChip(t)

	watcher, err := gpio.NewLineWatcher()
	assert.NoError(t, err, "unable to create LineWatcher")

	err = watcher.Add(c, line, gpio.RisingEdge, "testEventLineWait")
	assert.NoError(t, err, "unable to add event to the LineWatcher")

	event, err := watcher.Wait()
	assert.NoError(t, err, "unable to Wait on LineWatcher")
	assert.True(t, event.IsRising(), "trapped event is not of expected type")

	err = watcher.Close()
	assert.NoErrorf(t, err, "error while closing LineWatcher")
	c.Close()

	done <- struct{}{}
}
