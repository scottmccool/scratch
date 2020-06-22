package publishers

// Publish readings to PubNub
// TODO: This is a PoC, no error handling or cleanup is done at all

import (
	"os"

	pubnub "github.com/pubnub/go"

	"github.com/scottmccool/FBeacon/readings"
)

const pnChannelName = "FBeacon"

func PublishPubNub(pnChan chan readings.FBeacon) {
	config := pubnub.NewConfig()
	config.PublishKey = "pub-c-c7313055-d589-4a18-8bc3-2bc0a21d3b20"
	config.SubscribeKey = "sub-c-bd5d7130-b4a7-11ea-afa6-debb908608d9"
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
			_, _, _ = res, status, err
			//fmt.Printf("Pubblished to [[%v]]  (Res: [[%v]]) (Status: [[%v]]) (Err: [[%v]])\n", pnChannelName, *res, status, err)
		}
	}
}
