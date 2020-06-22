package publishers

//  gcloud pubsub topics create fbeacon-raw
// gcloud iam service-accounts create fbeacon-testing --description="FBeacon gcloud testing" --display-name="FBeacon-testing"
//gcloud iam service-accounts keys create secrets/fbeacon-testing --iam-account fbeacon-testing@fbeacon.iam.gserviceaccount.com
// grant it permissions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/scottmccool/FBeacon/readings"
)

func PublishGCloud(observations []readings.FBeacon) error {
	projectID := "fbeacon"
	topicID := "projects/fbeacon/topics/fbeacon-raw"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	t := client.Topic(topicID)
	var results []*pubsub.PublishResult
	for o := range observations {
		j, _ := json.Marshal(observations[o])
		r := t.Publish(ctx, &pubsub.Message{
			Data: []byte(string(j)),
		})
		fmt.Println("** Published: ", string(j))
		results = append(results, r)
	}
	time.Sleep(5 * time.Second)
	for _, r := range results {
		id, err := r.Get(ctx) // block until we get an ack
		if err != nil {
			fmt.Println("Error waiting on publish result: ", err)
			return err
			// TODO: Handle error.
		}
		fmt.Printf("Published a message with a message ID: %s\n", id)
	}
	return nil
}
