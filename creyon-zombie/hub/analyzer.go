package hub

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
	var readings []FBeacon
forLoop:
	for {
		select {
		case obs := <-Rawc:
			readings = append(readings, obs)
		default:
			break forLoop
		}
	}

	// Analyze them (nothing implemented!)

	// Publish them
	for reading := range readings {
		Analyzedc <- readings[reading]
	}
}
