# chardevgpio - How to test ?

During development, the library is tested using the Linux kernel module **gpio-mockup** on an x86_64 environment.

For real world tests on a Raspberry, see command line utilities provided under cmd directory.

## Pre-requisites

These steps required root access on the test machine:

1. Mount **debugfs** filesystem: ```sudo mount -t debugfs none /sys/kernel/debug```
2. Load **gpio-mockup** module: ```sudo modprobe gpio-mockup gpio_mockup_ranges=0,10 gpio_mockup_named_lines=1```

Finally, check the user that will be used for testing has read/write access on ```/dev/gpiochip0```.

## Running tests

It's easy as ```make test```
