package hub

// Publish analyzed readings to the cloud

import (
	"encoding/json"
	"fmt"
)

// Publish Reads from chan Analyzedc and publishes off hub (well, to stdout)
func Publish() {
	for {
		select {
		case obs := <-Analyzedc:
			//			fmt.Println(obs)
			j, err := json.Marshal(obs)
			if err != nil {
				fmt.Println("Cannot publish:", obs)
			} else {
				fmt.Println(string(j))
			}
		default:
			return
		}
	}
}
