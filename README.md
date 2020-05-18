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

Lines information can be requested from the chip at any moment as long as it is open.

```go
li, _ := chip.LineInfo(0)
line0Name := li.Name()
```

### HandleRequest

An HandleRequest is mandatory to setup a request an input line or an output line from the chip. The request should at minimum define the offsets of requested lines and the communication direction.

```go
lineIn_0, _    := chip.NewHandleRequest([]int{0}, gpio.HandleRequestInput)
lineOut_8_9, _ := chip.NewHandleRequest([]int{8, 9}, gpio.HandleRequestOutput)
```

Consumer name or default values can be set on the HandleRequest:

```go
lineIn_0.WithConsumer("myapp")
lineOut_8_9.WithDefaults([]int{1, 1})
```

Requesting lines to the chip is done passing a reference to the HandleRequest to ```Chip.RequestLines()```:

```go
c.RequestLines(lineIn_0)
c.RequestLines(lineOut_8_9) 
```

If no errors, the returned HandleRequest can then be used to read from or write to the lines:

```go
val0 := lineIn_0.Read()
lineOut_8_9.Write(0, 0)
```

Note that ```HandleRequest.Read()``` returns 3 values:

* first one is the read value for the first line managed by the HandleRequest. It is useful when working on a request with only one line.
* second one is an array containing read values for all lines managed by the HandleRequest
* last one is the error if any

### LineWatcher

Event on an input line can be trapped using a LineWatcher:

```go
watcher, _ := gpio.NewLineWatcher()
```

Add to it lines and events to be watched:

```go
watcher.Add(c, 0, gpio.RisingEdge, "wait rising edge on line 0")
watcher.Add(c, 1, gpio.FallingEdge, "wait falling edge on line 1")
watcher.Add(c, 2, gpio.BothEdges, "wait state change on line 2")
```

The watcher can then run indefinitely and call a handler function for each event trapped:

```go
func myFuncHandler(evd Event) { ... }
watcher.WaitForEver(myFuncHandler)
```

For simpler cases, the watcher can be used to block waiting for the first event occurrence:

```go
event, _ := watcher.Wait()
```

## Tests

During development, the library is tested using the Linux kernel module **gpio-mockup** on an x86_64 environment.

```
> ./setup-test.sh   # to be run once, prompts for sudo password
> make test
```

For real world tests on a Raspberry, see command line utilities provided under cmd directory.
