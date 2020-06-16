package hub

import (
	"fmt"
	"time"
)

const scanWaitT = 10 * time.Second   // Pause this long after getting a packet, save bt radio time
const analyzeBatchSize = 1           // Get this many readings before publishing
const pubFrequency = 1 * time.Second // Publish off hub this frequently

// Start activates a BLE scanner which will grab fujitsu packets from BLE beacons
func Start() {
	rawc := make(chan string, 10)       // readings from scanner to analyzer; buffered just in case analysis takes 10*ScanWait_t sec
	analyzedc := make(chan string, 100) // analyzer to publisher; buffered as publish may lag if we are conserving radio

	// Start a routine to publish analyzed readings
	go func() {
		for now := range time.Tick(pubFrequency) {
			_ = now
			publish(analyzedc)
		}
	}()

	//anonymous routine to handle data pipeline bits of analysis
	//Once you're ready to analyze raw data change this to go run analyze(rawc)
	//and implement analyze(rawc chan string)
	go func() {
		for reading := range rawc {
			// insert fancy analysis here
			analyzedc <- reading
		}
	}()

	// main thread scans for packets which it'll pass to analyzer
	// TODO: Gotta be a more elegant way to do this!
	for {
		start := time.Now()
		rawc <- scan() // Read a BTLE ad and push to analyzer, sleep
		elapsed := time.Since(start)
		sleept := scanWaitT - elapsed
		if sleept > 0 {
			//			fmt.Println("Scan took ", elapsed, ".  Sleeping ", sleept.Seconds(), " for next scan")
			time.Sleep(sleept)
		}
	}
	select {} // block forever
}

// Scans for a env reading packet up to X seconds; return first match
func scan() string {
	return "{\"mock\": true, \"timestamp\": \"" + time.Now().String() + "\"}"
}

func publish(analyzedc chan string) {
	reading := <-analyzedc
	fmt.Println("Publishing: ", reading)
}
