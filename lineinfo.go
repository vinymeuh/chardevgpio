// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// +build linux

package chardevgpio

/** LineInfo structure is defined in linux.go **/

// IsOutput returns true if the line is configured as an output.
func (li LineInfo) IsOutput() bool {
	return li.Flags&gpioLineFlagIsOut == gpioLineFlagIsOut
}

// IsInput returns true if the line is configured as an input.
func (li LineInfo) IsInput() bool {
	return !li.IsOutput()
}

// IsActiveLow returns true if the line is configured as active low.
func (li LineInfo) IsActiveLow() bool {
	return li.Flags&gpioLineFlagActiveLow == gpioLineFlagActiveLow
}

// IsActiveHigh returns true if the line is configured as active high.
func (li LineInfo) IsActiveHigh() bool {
	return !(li.IsActiveLow())
}

// IsOpenDrain returns true if the line is configured as open drain.
func (li LineInfo) IsOpenDrain() bool {
	return li.Flags&gpioLineFlagOpenDrain == gpioLineFlagOpenDrain
}

// IsOpenSource returns true if the line is configured as open source.
func (li LineInfo) IsOpenSource() bool {
	return li.Flags&gpioLineFlagOpenSource == gpioLineFlagOpenSource
}

// IsKernel returns true if the line is configured as kernel.
func (li LineInfo) IsKernel() bool {
	return li.Flags&gpioLineFlagKernel == gpioLineFlagKernel
}
