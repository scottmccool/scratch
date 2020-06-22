package analyzers

import (
	"github.com/scottmccool/FBeacon/hub/sniffers"
	"github.com/scottmccool/FBeacon/readings"
)

// Analyzes batches of readings
// May perform event detection (occupancy) or filtering
// TODO: At least dedupe!
// For now just pass through, we will practice batching in publish
const analyzeMinBatchSize = 5

var AnalyzedReadings = make(chan readings.FBeacon, 1000)

// Analyze read from chan Rawc once analyzeMinBatchSize readings are available, filters (not done) and passes to publisher channel
func Analyze() (analyzed int, err error) {
	//log.Printf("Beginning analysis (Rawc: %v)\n", len(Rawc))
	if len(sniffers.SniffedObservations) < analyzeMinBatchSize {
		return 0, nil // Not enough readings to analyze, try again later
	}

	// Read all available data
	var obs []readings.FBeacon
forLoop:
	for {
		select {
		case o := <-sniffers.SniffedObservations:
			obs = append(obs, o)
		default:
			break forLoop
		}
	}

	// Analyze them (nothing implemented!)

	// Publish them
	for o := range obs {
		AnalyzedReadings <- obs[o]
	}

	return len(obs), nil
}
