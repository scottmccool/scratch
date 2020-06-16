Reimplement creyon in go?

hub.go :
  scanner -- Every X seconds, read bt and decode fujitsu data; publish to (buffered) channel
  recorder -- Read scanner results from channel until payload size is reached; publish to (unbuffered) channel
  launcher -- Read a payload from recorder (set of records) and publish it to pubnub environmental data channel

analysis.go
  Subscribe to environmental data, update datadog chart
