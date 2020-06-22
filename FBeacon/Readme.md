# FBeacon scanner

## What is it?

A simple data pipeline intended to run on a Raspberry Pi with BLE.

It looks for Fujitsu beacons containing environmental data, translates the data, filters it, and publishes readings.

Uses goroutines for concurrency, supports multiple scanners, analyzers, and publishers

## Usage

```bash
go install

# consider adding 2> /dev/null to avoid weird data lines
sudo -e FBeacon
```

or

```bash
sudo go run main.go
```

from this directory.

## Background

Reads environmental data from a Fujitsu BLE tag.

hub.go:
  Manage three main threads of execution:
    *fuji_tag_scanner.go -- Scans btle, finds fuji beacons and extracts the data into a simple struct, passes along to analyzer
    * analyzer.go -- Reads a batch of readings, filters or modifies them, and sends them on to publisher
    * publisher.go -- For now, this just prints some interesting stuff about the tag readings to stdout

Channels are used to form a basic data pipeline
  ~~btle radio waves~~ --> fuji_tag_scanner --> analyzer --> publisher

Batching and cost optimization (battery, connectivity, cloud) not yet implemented, but would be done by adjusting batch sizes and timer frequencies.
