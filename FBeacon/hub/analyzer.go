package hub

import (
	"github.com/scottmccool/FBeacon/readings"
)

// Analyzes batches of readings
// May perform event detection (occupancy) or filtering
// For now just pass through, we will practice batching in publish

// Analyze read from chan Rawc once analyzeMinBatchSize readings are available, filters (not done) and passes to publisher channel
func Analyze() {
	//log.Printf("Beginning analysis (Rawc: %v)\n", len(Rawc))
	if len(Rawc) < analyzeMinBatchSize {
		return // Not enough readings to analyze, try again later
	}

	// Read all available data
	var obs []readings.FBeacon
forLoop:
	for {
		select {
		case o := <-Rawc:
			obs = append(obs, o)
		default:
			break forLoop
		}
	}

	// Analyze them (nothing implemented!)

	// Publish them
	for o := range obs {
		Analyzedc <- obs[o]
	}
}
