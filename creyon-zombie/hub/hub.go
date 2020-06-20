package hub

import (
	"fmt"
	"log"
	"time"
)

// FBeacon - a Fujitsu tag
type FBeacon struct {
	temp         float32
	xAcc         float32
	yAcc         float32
	zAcc         float32
	addr         string
	txPowerLevel int
	rawMfrData   []byte
	timestamp    time.Time
}

const scanWaitT = 1 * time.Second    // Pause this long after getting a packet, save bt radio time
const analyzeBatchSize = 2           // Get this many readings before publishing
const pubFrequency = 5 * time.Second // Publish off hub this frequently

// BT handler needs to communicate; simplest to export the channels
var Rawc = make(chan FBeacon, 500)
var Analyzedc = make(chan FBeacon, 500)

// Start activates a BLE scanner which will grab fujitsu packets from BLE beacons
func Start() {
	log.Printf("Scanning for Fuji tags every %v, publishing to cloud every %v after analysis in batches of %v\n", scanWaitT, analyzeBatchSize, pubFrequency)
	//rawc := make(chan string, 50)      // readings from scanner to analyzer; buffered just in case analysis takes 10*ScanWait_t sec
	//analyzedc := make(chan string, 50) // analyzer to publisher; buffered as publish may lag if we are conserving radio

	// Start a routine to publish analyzed readings
	go func() {
		for now := range time.Tick(pubFrequency) {
			_ = now
			publish()
		}
	}()

	// And another to read from scanner and do on-hub processing before publish
	go func() {
		for now := range time.Tick(1 * time.Second) {
			_ = now
			analyze()
		}
	}()

	// main thread scans for packets which it'll pass to analyzer
	// TODO: Gotta be a more elegant way to do this!
	go func() {
		for {
			start := time.Now()
			scan() // No real timer logic implemented yet
			elapsed := time.Since(start)
			sleept := scanWaitT - elapsed
			if sleept > 0 {
				fmt.Println("Scan took ", elapsed, ".  Sleeping ", sleept.Seconds(), " for next scan")
				time.Sleep(sleept)
			}
		}
	}()

	select {} // block forever
}

// Scans for a env reading packet up to X seconds; return first match
func scan() {
	//return "{\"mock\": true, \"timestamp\": \"" + time.Now().String() + "\"}"
	ScanFuji()
}

// Analyzes batches of readings
// May perform event detection (occupancy) or filtering
// For now just pass through, we will practice batching in publish
func analyze() { //rawc chan string, analyzedc chan string) {
	if len(Rawc) < analyzeBatchSize {
		return // Not enough readings to analyze, try again later
	}

	// Read everything a batch
	var readings []FBeacon
	for len(readings) < analyzeBatchSize {
		select {
		case reading := <-Rawc:
			readings = append(readings, reading)
			//		default:
			//			time.Sleep(1 * time.Second)
		}
	}

	// Analyze them!

	// Publish them
	for reading := range readings {
		Analyzedc <- readings[reading]
	}
}

// Publish analyzed readings to the cloud
// For now, just to stdout.
// Non-blocking, fired by timer should read all available and publish
func publish() { //analyzedc chan string) {
	for {
		select {
		case reading := <-Analyzedc:
			fmt.Println("Publishing: ", reading)
		default:
			return
		}
	}
}
