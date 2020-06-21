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

	gcloud_err := publishToGcloud(observations)
	stdout_err := publishToStdOut(observations)

	if gcloud_err != nil {
		return len(observations), gcloud_err
	}
	if stdout_err != nil {
		return len(observations), stdout_err
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
