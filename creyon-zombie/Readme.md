# Creyon

Reimplement creyon in go?

Reads environmental data from a Fujitsu BLE tag.

hub.go:
  Manage three main threads of execution:
    *fuji_tag_scanner.go -- Scans btle, finds fuji beacons and extracts the data into a simple struct, passes along to analyzer
    * analyzer.go -- Reads a batch of readings, filters or modifies them, and sends them on to publisher
    * publisher.go -- For now, this just prints some interesting stuff about the tag readings to stdout

Channels are used to form a basic data pipeline
  ~~btle radio waves~~ --> fuji_tag_scanner --> analyzer --> publisher

Batching and cost optimization (battery, connectivity, cloud) not yet implemented, but would be done by adjusting batch sizes and timer frequencies.
