package hub

import (
	"fmt"
)

// Publish analyzed readings to the cloud
// For now, just to stdout.
// publish() Reads from chan Analyzedc and publishes off hub (well, to stdout)
func Publish() {
	for {
		select {
		case obs := <-Analyzedc:
			fmt.Println(obs)
		default:
			return
		}
	}
}
