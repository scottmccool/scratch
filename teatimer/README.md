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
go get -u github.com/spf13/cobra/cobra
asdf reshim golang
//go mod init tea-timer-app
cobra init --pkg-name tea-timer-app -a "Scott McCool"
cobra add start
go run main.go start
//go install tea-timer-app

```


## References
[1] https://github.com/spf13/cobra
[2] https://towardsdatascience.com/how-to-create-a-cli-in-golang-with-cobra-d729641c7177
