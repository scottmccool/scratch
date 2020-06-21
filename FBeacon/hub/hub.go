package hub

import (
	"github.com/scottmccool/FBeacon/readings"
	"log"
	"time"
)

// How frequently we sniff and publish sensor packets
const analyzeMinBatchSize = 1        // Get this many readings from sniffer before analyzing (batch size for on hub analysis, can use to batch publishes too)
const pubFrequency = 1 * time.Second // Publish everything off hub this frequently (analysis will write to the analyzed channel this reads from)

// Rawc channel for sniffer to analzyer communications
var Rawc = make(chan readings.FBeacon, 1000)

// Analyzedc channel for analyzer to publisher communications
var Analyzedc = make(chan readings.FBeacon, 50)

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
		for now := range time.Tick(1 * time.Second) {
			_ = now
			Analyze()
		}
	}()

	// main thread scans for packets which it'll pass to analyzer
	go func() {
		for {
			ScanFuji() // Activate bluetooth for scanT time
		}
	}()
	select {} // block forever
}
