package hub

import (
	"log"
	"time"

	"github.com/scottmccool/FBeacon/hub/sniffers"
)

// How frequently we sniff and publish sensor packets
const analyzeMinBatchSize = 2         // Get this many readings from sniffer before analyzing (batch size for on hub analysis, can use to batch publishes too)
const pubFrequency = 10 * time.Second // Publish everything off hub this frequently (analysis will write to the analyzed channel this reads from)
const pubMinBatchSize = 10            // Wait to accumulate MinBatchSize before publishing
const analyzeFrequency = 5 * time.Second

// Start Manage three worker routines (scanner, analyzer, publisher)
func Start() {
	log.Printf("Scanning for Fuji sensor tags; analyzing in batches of %v and publishing (printing) every %v\n", analyzeMinBatchSize, pubFrequency)

	// Start a routine to publish analyzed readings
	go func() {
		for now := range time.Tick(pubFrequency) {
			_ = now
			Publish()
		}
	}()

	// And another to read from scanner and do on-hub processing before publish, once a second
	go func() {
		for now := range time.Tick(analyzeFrequency) {
			_ = now
			Analyze()
		}
	}()

	sniffers.ScanFuji() // Run forever.
}
