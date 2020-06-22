package publishers

// Publish readings to PubNub
// PoC

import (
	"fmt"
	"os"

	pubnub "github.com/pubnub/go"

	"github.com/scottmccool/FBeacon/readings"
)

const pnChannelName = "FBeacon"

func PublishPubNub(pnChan chan readings.FBeacon) {
	config := pubnub.NewConfig()
	config.PublishKey = "pub-c-1dadcdf6-f319-4eac-9d81-32ac076c6791"
	config.SubscribeKey = "sub-c-d8802764-b439-11ea-afa6-debb908608d9"
	hn, err := os.Hostname()
	if err != nil {
		config.UUID = "unknown-hub"
	} else {
		config.UUID = hn
	}

	pn := pubnub.NewPubNub(config)

	for {
		select {
		case obs := <-pnChan:
			res, status, err := pn.Publish().
				Channel(pnChannelName).
				Message(obs).
				Execute()
			fmt.Printf("Pubblished to [[%v]]  (Res: [[%v]]) (Status: [[%v]]) (Err: [[%v]])\n", pnChannelName, *res, status, err)
		}
	}
}
