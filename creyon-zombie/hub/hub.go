package hub

import "fmt"

// Start activates a BLE scanner which will grab fujitsu packets from BLE beacons
func Start() {
	fmt.Println("Starting publisher thread for analyzed events")
	fmt.Println("Starting analyzer thread for scanned events")
	fmt.Println("Starting scanner thread to do the bluetooth things!")
}

// Scans for a env reading packet every X seconds
func scan(interval int) {
	// every interval seconds, scan bt for a matching packet
}

// Accumulates env readings from scanner
// Filter env readings and publish them
func analyze() {
	// Accumulate readings from scanner, find interesting events
	// For a batch of N readings, see if occupied changes from true or false and fire that as a different event
	// At N readings, publish to firehose
}

// Publishes a batch of analyzed readings up to the cloud for archival and broader analysis
func publish() {
	// Read from firehose and events channels
	// Publish to corresponding pubnub channels
}
