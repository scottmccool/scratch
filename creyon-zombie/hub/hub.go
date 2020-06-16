package hub

import (
	"fmt"
	"time"
)

const scanWaitT = 5 * time.Second      // Pause this long after getting a packet, save bt radio time
const analyzeBatchSize = 2 // Get this many readings before publishing

// Start activates a BLE scanner which will grab fujitsu packets from BLE beacons
func Start() {
	rawc := make(chan string, 10)       // readings from scanner to analyzer; buffered just in case analysis takes 10*ScanWait_t sec
	analyzedc := make(chan string, 100) // analyzer to publisher; buffered as publish may lag if we are conserving radio

	//anonymous routine to handle data pipeline bits of analysis
	go func() {
		for reading := range rawc {
			// insert fancy analysis here
			analyzedc <- reading
		}
	}()

	go func() {
		for reading := range analyzedc {
			fmt.Println(reading)
		}
	}()

	// main thread scans for packets which it'll pass to analyzer
	// TODO: Gotta be a more elegant way to do this!
	for {
		start := time.Now()
		rawc <- scan() // Read a BTLE ad and push to analyzer, sleep
		elapsed := time.Since(start)
		sleept := scanWaitT - elapsed
		fmt.Println("Scan took ", elapsed, ".  Sleeping ", sleept.Seconds(), " for next scan")
		time.Sleep(sleept)
	}
	select {} // block forever
}

// Scans for a env reading packet up to X seconds; return first match
func scan() string {
	return "Mock reading:" + time.Now().String()
}

// Accumulates env readings from scanner
// Filter env readings and publish them
func analyze(rawc, analyzedc chan string) {
	//var readings[analyzeBatchSize] string
	// Loop forever, accumulating a batch of readings and publishing
	// Insert event detection logic here
	for {
		//fmt.Println("Skiping analysis and passing to publish")
		reading := <-rawc
		analyzedc <- reading
	}
}

// Accumulate readings from scanner, find interesting events
// For a batch of N readings, see if occupied changes from true or false and fire that as a different event
// At N readings, publish to firehose

// Publishes a batch of analyzed readings up to the cloud for archival and broader analysis
func publish(c chan string) {
	// Read from firehose and events channels
	// Publish to corresponding pubnub channels
	reading := <-c
	fmt.Println(reading)
}
