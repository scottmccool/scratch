# TeaTimer

## Overview
A test project to learn go and the Coba CLI framework

## Specifications
* CLI application to start, show, and cancel a timer.
* Should be pub-nub'able in the future

Example Usage:
```bash
 $ ./teatimer help
----
Manage a tea timer!

Usage:
  ./teatimer [start|show|cancel] <seconds>

Commands:
start: Start a timer (error if one is active)
show: Show seconds remaining (or elapsed since expiration; error if none started ever)
cancel: Cancel any active timer

```

Use [Cobra](https://github.com/spf13/cobra) CLI framework:
```bash
$ go get -u github.com/spf13/cobra/cobra
$ asdf reshim golang
$ cobra --init --pkg-name teatimer -a "Scott McCool" --viper
```
