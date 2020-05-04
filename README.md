# chardevgpio

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/vinymeuh/chardevgpio.svg)](https://github.com/vinymeuh/chardevgpio/releases/latest)
[![GoDoc](https://godoc.org/github.com/vinymeuh/chardevgpio?status.svg)](https://godoc.org/github.com/vinymeuh/chardevgpio)
[![Build Status](https://travis-ci.org/vinymeuh/chardevgpio.svg?branch=master)](https://travis-ci.org/vinymeuh/chardevgpio)
[![codecov](https://codecov.io/gh/vinymeuh/chardevgpio/branch/master/graph/badge.svg)](https://codecov.io/gh/vinymeuh/chardevgpio)
[![Go Report Card](https://goreportcard.com/badge/github.com/vinymeuh/chardevgpio)](https://goreportcard.com/report/github.com/vinymeuh/chardevgpio)

**chardevgio** is a pure Go library for access the Linux GPIO character device user API.

## Usage

```go
import gpio "github.com/vinymeuh/chardevgpio"
```

In following examples, error handling will be ommited.

### Chip

A Chip object must be open before to request lines or lines informations from a GPIO chip.

```go
chip, _ := gpio.NewChip("/dev/gpiochip0")
...
defer chip.Close()
```

Closing a chip does not invalidate any previously requested lines that can still be used.

### LineInfo

Lines information can be requested from the chip as long as it is open.

```go
li, _ := chip.LineInfo(0)
line0Name := li.Name()
```
