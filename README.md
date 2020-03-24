# chardevgpio

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/vinymeuh/chardevgpio?status.svg)](https://godoc.org/github.com/vinymeuh/chardevgpio)
[![Build Status](https://travis-ci.org/vinymeuh/chardevgpio.svg?branch=master)](https://travis-ci.org/vinymeuh/chardevgpio)
[![Go Report Card](https://goreportcard.com/badge/github.com/vinymeuh/chardevgpio)](https://goreportcard.com/report/github.com/vinymeuh/chardevgpio)

**chardevgio** is a pure Go library for access the Linux GPIO character device user API.

## Usage

**chardevgio** provides two level of API:

* **GPIOx** structures which mimics directly the Linux C code and must be used directly with IOCTL syscalls
* A thin level of abstraction over **GPIOx** structures to hide IOCTL syscalls and simplify use
