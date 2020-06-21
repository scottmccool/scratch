package hub

import (
	"github.com/scottmccool/FBeacon/readings"
)

// publishToGcloud Write a batch of records to pubsub topic
func publishToGcloud(observations []readings.FBeacon) error {
	// Write a slice of
	// {"Timestamp":"2020-06-20T22:18:49.450340744-07:00","BtData":{"Addr":"E2:94:B4:AF:93:13","TxPowerLevel":0,"Rssi":-79,"RawMfrData":"59000100030003004401270010009308"},"Measurements":{"Temp":71.54679,"XAcc":0.019042969,"YAcc":0.0078125,"ZAcc":1.0717773}}
	// to pubsub for drain to bigquery

	return nil
}
