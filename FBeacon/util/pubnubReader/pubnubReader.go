package main

import (
	"fmt"
	"os"

	pubnub "github.com/pubnub/go"
)

func main() {
	// Initialize pn
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

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case signal := <-listener.Signal:
				//Channel
				fmt.Println(signal.Channel)
				//Subscription
				fmt.Println(signal.Subscription)
				//Payload
				fmt.Println(signal.Message)
				//Publisher ID
				fmt.Println(signal.Publisher)
				//Timetoken
				fmt.Println(signal.Timetoken)
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNDisconnectedCategory:
					// this is the expected category for an unsubscribe. This means there
					// was no error in unsubscribing from everything
				case pubnub.PNConnectedCategory:
					// this is expected for a subscribe, this means there is no error or issue whatsoever
				case pubnub.PNReconnectedCategory:
					// this usually occurs if subscribe temporarily fails but reconnects. This means
					// there was an error but there is no longer any issue
				case pubnub.PNAccessDeniedCategory:
					// this means that PAM does allow this client to subscribe to this
					// channel and channel group configuration. This is another explicit error
				}
			case message := <-listener.Message:
				//Channel
				fmt.Println(message.Channel)
				//Subscription
				fmt.Println(message.Subscription)
				//Payload
				fmt.Println(message.Message)
				//Publisher ID
				fmt.Println(message.Publisher)
				//Timetoken
				fmt.Println(message.Timetoken)
			case presence := <-listener.Presence:
				fmt.Println(presence.Event)
				//Channel
				fmt.Println(presence.Channel)
				//Subscription
				fmt.Println(presence.Subscription)
				//Timetoken
				fmt.Println(presence.Timetoken)
				//Occupancy
				fmt.Println(presence.Occupancy)
			case membershipEvent := <-listener.MembershipEvent:
				fmt.Println(fmt.Sprintf("membershipEvent.Channel: %s", membershipEvent.Channel))
				fmt.Println(fmt.Sprintf("membershipEvent.SubscribedChannel: %s", membershipEvent.SubscribedChannel))
				fmt.Println(fmt.Sprintf("membershipEvent.Event: %s", membershipEvent.Event))
				fmt.Println(fmt.Sprintf("membershipEvent.Description: %s", membershipEvent.Description))
				fmt.Println(fmt.Sprintf("membershipEvent.Timestamp: %s", membershipEvent.Timestamp))
				fmt.Println(fmt.Sprintf("membershipEvent.Custom: %v", membershipEvent.Custom))
			case messageActionsEvent := <-listener.MessageActionsEvent:
				fmt.Println(fmt.Sprintf("messageActionsEvent.Channel: %s", messageActionsEvent.Channel))
				fmt.Println(fmt.Sprintf("messageActionsEvent.SubscribedChannel: %s", messageActionsEvent.SubscribedChannel))
				fmt.Println(fmt.Sprintf("messageActionsEvent.Event: %s", messageActionsEvent.Event))
				fmt.Println(fmt.Sprintf("messageActionsEvent.Data.ActionType: %s", messageActionsEvent.Data.ActionType))
				fmt.Println(fmt.Sprintf("messageActionsEvent.Data.ActionValue: %s", messageActionsEvent.Data.ActionValue))
				fmt.Println(fmt.Sprintf("messageActionsEvent.Data.ActionTimetoken: %s", messageActionsEvent.Data.ActionTimetoken))
				fmt.Println(fmt.Sprintf("messageActionsEvent.Data.MessageTimetoken: %s", messageActionsEvent.Data.MessageTimetoken))
			}
		}
	}()

	// Add subscriber
	// pn.Subscribe().
	// 	Channels([]string{"FBeacon"}). // subscribe to channels
	// 	Execute()

	res, status, err := pn.Time().Execute()
	fmt.Println(res, status, err)
	hres, hstatus, herr := pn.History().
		Channel("FBeacon"). // where to fetch history from
		Count(10).          // how many items to fetch
		Execute()

	fmt.Println(hres, hstatus, herr)
	//	select {}
}
