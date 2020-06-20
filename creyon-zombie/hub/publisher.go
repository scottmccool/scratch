package hub

// Publish analyzed readings to the cloud

import (
	"fmt"
)

// Publish() Reads from chan Analyzedc and publishes off hub (well, to stdout)
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
