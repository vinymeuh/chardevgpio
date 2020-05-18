#!/usr/bin/env bash
#
# Mount debugfs and load gpio-mockup kernel module.
# 

exe() { echo "$@" ; "$@" ; }

if [ $(grep -c /sys/kernel/debug /proc/mounts) -eq 0 ]; then
    exe sudo mount -t debugfs none /sys/kernel/debug
fi

if [ $(lsmod | grep -c gpio_mockup) -eq 0 ]; then
    exe modprobe gpio-mockup gpio_mockup_ranges=0,10 gpio_mockup_named_lines=1
fi

exe sudo chmod a+rx /sys/kernel/debug
exe sudo chmod u+rw /sys/kernel/debug/gpio-mockup/gpiochip0/*

user=$(whoami)

exe sudo chown $user /dev/gpiochip0
exe sudo chown -R $user /sys/kernel/debug/gpio*
