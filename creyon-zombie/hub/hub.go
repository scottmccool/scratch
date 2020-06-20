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
	rssi         int
	rawMfrData   string
	timestamp    time.Time
}

// stringer for observations
func (o FBeacon) String() string {
	// metadata, Temp:value, Acc:value
	return fmt.Sprintf("(%v):[%v](%v), Temp: %v, Acc: (%v, %v, %v)", o.timestamp.Format(time.Stamp), o.addr, o.rssi, o.temp, o.xAcc, o.yAcc, o.zAcc)
	//	return fmt.Sprintf("%v (%v years)", p.Name, p.Age)
}

// How frequently we sniff and publish sensor packets
const analyzeMinBatchSize = 1        // Get this many readings from sniffer before analyzing (batch size for on hub analysis, can use to batch publishes too)
const pubFrequency = 1 * time.Second // Publish everything off hub this frequently (analysis will write to the analyzed channel this reads from)

// Rawc channel for sniffer to analzyer communications
var Rawc = make(chan FBeacon, 1000)

// Analyzedc channel for analyzer to publisher communications
var Analyzedc = make(chan FBeacon, 50)

// HubStart Manage three worker routines (scanner, analyzer, publisher)
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
	// TODO: Gotta be a more elegant way to do this!
	// For now assume scanner will exit and we want to wait at least scanWaitT between BLE radio activation
	go func() {
		for {
			ScanFuji() // Activate bluetooth for scanT time
		}
	}()
	select {} // block forever
}
