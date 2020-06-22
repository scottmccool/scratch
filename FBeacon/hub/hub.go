package hub

import (
	"log"
	"time"

	"github.com/scottmccool/FBeacon/hub/analyzers"
	"github.com/scottmccool/FBeacon/hub/publishers"
	"github.com/scottmccool/FBeacon/hub/sniffers"
)

// How frequently we sniff and publish sensor packets
const pubFrequency = 10 * time.Second // Publish everything off hub this frequently (analysis will write to the analyzed channel this reads from)
const analyzeFrequency = 1 * time.Second

// Start Manage three worker routines (scanner, analyzer, publisher)
func Start() {
	log.Printf("Starting collector software, looking for Fujitsu BLE beacons.")

	// Start a routine to publish analyzed readings
	go func() {
		for now := range time.Tick(pubFrequency) {
			_ = now
			_, err := publishers.Publish()
			if err != nil {
				//fmt.Println("Swallowing publihs error(s): ", err)
			} else {
				//fmt.Println("Published ", published, " observations")
			}
		}
	}()

	// And another to read from scanner and do on-hub processing before publish, once a second
	go func() {
		for now := range time.Tick(analyzeFrequency) {
			_ = now
			_, err := analyzers.Analyze()
			if err != nil {
				//fmt.Println("Swallowing analysis error(s): ", err)
			} else {
				//fmt.Println("Analyzed ", analyzed, " observations")
			}

		}
	}()

	sniffers.ScanFuji() // Run forever.
}
