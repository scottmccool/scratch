package publishers

// Publish analyzed readings to the cloud

import (
	"github.com/scottmccool/FBeacon/hub/analyzers"
	"github.com/scottmccool/FBeacon/readings"
)

const pubMinBatchSize = 10

// TODO: Defer close until any child publishers are done

// Publish Reads from chan Analyzedc and publishes off hub (well, to a set of publishers also managed here
// This needs refactored for real world use, stop recreating clients to external services so often
func Publish() (published int, err error) {
	if len(analyzers.AnalyzedReadings) < pubMinBatchSize {
		return 0, nil
	}
	// Spin up publishers
	pubNubChan := make(chan readings.FBeacon, len(analyzers.AnalyzedReadings))
	loggerChan := make(chan readings.FBeacon, len(analyzers.AnalyzedReadings))
	go PublishPubNub(pubNubChan)
	go PublishLog(loggerChan)

	// Read all available data
	var obs []readings.FBeacon
forLoop:
	for {
		select {
		case o := <-analyzers.AnalyzedReadings:
			pubNubChan <- o
			loggerChan <- o
		default:
			break forLoop
		}
	}

	return len(obs), nil
}
