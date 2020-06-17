package hub

import (
	"fmt"
	"time"
)

const scanWaitT = 1 * time.Second    // Pause this long after getting a packet, save bt radio time
const analyzeBatchSize = 3           // Get this many readings before publishing
const pubFrequency = 5 * time.Second // Publish off hub this frequently

// Start activates a BLE scanner which will grab fujitsu packets from BLE beacons
func Start() {
	fmt.Println("Scanning every ", scanWaitT, ", analyzing batches of ", analyzeBatchSize, " and publishing every ", pubFrequency)
	rawc := make(chan string, 50)      // readings from scanner to analyzer; buffered just in case analysis takes 10*ScanWait_t sec
	analyzedc := make(chan string, 50) // analyzer to publisher; buffered as publish may lag if we are conserving radio

	// Start a routine to publish analyzed readings
	go func() {
		for now := range time.Tick(pubFrequency) {
			_ = now
			publish(analyzedc)
		}
	}()

	// And another to read from scanner and do on-hub processing before publish
	go analyze(rawc, analyzedc)

	// main thread scans for packets which it'll pass to analyzer
	// TODO: Gotta be a more elegant way to do this!

	go func() {
		for {
			start := time.Now()
			rawc <- scan() // Read a BTLE ad and push to analyzer, sleep
			elapsed := time.Since(start)
			sleept := scanWaitT - elapsed
			if sleept > 0 {
				//fmt.Println("Scan took ", elapsed, ".  Sleeping ", sleept.Seconds(), " for next scan")
				time.Sleep(sleept)
			}
		}
	}()

	select {} // block forever
}

// Scans for a env reading packet up to X seconds; return first match
func scan() string {
	fmt.Println("faking it")
	return "{\"mock\": true, \"timestamp\": \"" + time.Now().String() + "\"}"
}

// Analyzes batches of readings
// May perform event detection (occupancy) or filtering
// Publishes any readings or signals one by one to analyzedc
func analyze(rawc chan string, analyzedc chan string) { // TODO This needs to loop until full.
	var readings [analyzeBatchSize]string
	for i := 0; i < analyzeBatchSize; i++ {
		fmt.Println("Reading it")
		select {
		case reading := <-rawc:
			readings[i] = reading
		default:
			time.Sleep(1 * time.Second)
		}
	}

	// we have a batch, just publish it
	for _, s := range readings {
		analyzedc <- s
	}
}

func publish(analyzedc chan string) {
	select {
	case reading := <-analyzedc:
		fmt.Println("Publishing: ", reading)
	default:
		return
	}
}
