package hub

import "fmt"

// Start activates a BLE scanner which will grab fujitsu packets from BLE beacons
func Start() {
	fmt.Println("Starting publisher thread")
	fmt.Println("Starting analyzer thread")
	fmt.Println("Starting scanner thread")
}

// Scans for a fuji env reading packet every X seconds
func scan(interval int) {

}

// Accumulates env readings from scanner
// Filter env readings and publish them
func analyze() {

}

// Publishes a batch of analyzed readings up to the cloud for archival and broader analysis
func publish() {

}

// grabBLEPacket activates the BLE interface and reads one broadcast packet from fujitsu sensors
// Returns a CreyonPayload
func grabBLEPacket() {
	// Grab a set of env readings
}
