package hub

// Publish analyzed readings to the cloud

import (
	"encoding/json"
	"fmt"

	"github.com/scottmccool/FBeacon/readings"
)

// Publish Reads from chan Analyzedc and publishes off hub (well, to stdout)
func Publish() (published int, err error) {
	if len(Analyzedc) < pubMinBatchSize {
		return 0, nil
	}
	var observations []readings.FBeacon
	for len(observations) < pubMinBatchSize {
		select {
		case obs := <-Analyzedc:
			observations = append(observations, obs)
		}
	}

	gcloudErr := publishToGcloud(observations)
	stdoutErr := publishToStdOut(observations)

	if gcloudErr != nil {
		return len(observations), gcloudErr
	}
	if stdoutErr != nil {
		return len(observations), stdoutErr
	}
	return len(observations), nil
}

func publishToStdOut(observations []readings.FBeacon) error {
	//	fmt.Println(observations)
	for o := range observations {
		s, err := json.Marshal(observations[o])
		if err != nil {
			return err
		}
		fmt.Println(string(s))
	}
	return nil
}
