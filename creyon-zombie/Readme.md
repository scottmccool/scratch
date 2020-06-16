# Creyon

Reimplement creyon in go?

hub.go :
  scan -- Every X seconds, read bt and decode fujitsu data; publish to (buffered) channel
  analyze -- Read scanner results from channel until payload size is reached; filter and pass to publisher
  publish -- Read a payload from recorder (set of records) and publish it to pubnub environmental data channel

Data "accumulates" in analyze.  It should be able to write unpublished readings to stdout on graceful shutdown; it should suicide at MAXUNPUBLISHEDREADINGS

analysis.go
  Subscribe to environmental data, update datadog chart or publish to slack or whatever
