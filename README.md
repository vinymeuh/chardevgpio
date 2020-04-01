# chardevgpio

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/vinymeuh/chardevgpio.svg)](https://github.com/vinymeuh/chardevgpio/releases/latest)
[![GoDoc](https://godoc.org/github.com/vinymeuh/chardevgpio?status.svg)](https://godoc.org/github.com/vinymeuh/chardevgpio)
[![Build Status](https://travis-ci.org/vinymeuh/chardevgpio.svg?branch=master)](https://travis-ci.org/vinymeuh/chardevgpio)
[![Go Report Card](https://goreportcard.com/badge/github.com/vinymeuh/chardevgpio)](https://goreportcard.com/report/github.com/vinymeuh/chardevgpio)

**chardevgio** is a pure Go library for access the Linux GPIO character device user API, providing two level of API:

* **GPIOx** structures which mimics the Linux C code and must be used directly with ioctl syscalls
* A thin level of abstraction over **GPIOx** structures to hide ioctl syscalls and simplify use

We describe only the use of the "higher" level API. See [chardevgpio.go](./chardevgpio.go) for examples of how to use the **GPIOx** structures.

In following examples, error handling will be ommited.

## Usage

```go
import gpio "github.com/vinymeuh/chardevgpio"
```

### Chip

A Chip object must be open before to request lines or lines informations from a GPIO chip.

```go
chip, _ := gpio.Open("/dev/gpiochip0")
...
defer chip.Close()
```

Once opened, chip can be used to access lines informations or to request lines.

Closing a chip does not invalidate any previously requested lines that can still be used.
