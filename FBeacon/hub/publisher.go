package hub

// Publish analyzed readings to the cloud

import (
	"github.com/scottmccool/FBeacon/hub/publishers"
	"github.com/scottmccool/FBeacon/readings"
)

// Publish Reads from chan Analyzedc and publishes off hub (well, to a set of publishers also managed here
// This needs refactored for real world use, stop recreating clients to external services so often
func Publish() (published int, err error) {
	if len(AnalyzedReadings) < pubMinBatchSize {
		return 0, nil
	}
	// Spin up publishers
	pnChan := make(chan readings.FBeacon, len(AnalyzedReadings))
	loggerChan := make(chan readings.FBeacon, len(AnalyzedReadings))
	go publishers.PublishPubNub(pnChan)
	go publishers.PublishLog(loggerChan)

	var observations []readings.FBeacon
	for len(observations) < pubMinBatchSize {
		select {
		case obs := <-AnalyzedReadings:
			observations = append(observations, obs)
		}
	}

	for _, obs := range observations {
		pnChan <- obs
		loggerChan <- obs
	}

	return len(observations), nil
}
