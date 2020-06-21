package readings

// Our sensor data reading object(s)

import (
	"fmt"
	"time"
)

// FBeacon a fujitsu tag measuring temp, x_acc, y_acc, z_acc along with bluetooth metadata
type FBeacon struct {
	Timestamp    time.Time
	BtData       btData
	Measurements measurements
}

type measurements struct {
	Temp float32
	XAcc float32
	YAcc float32
	ZAcc float32
}

type btData struct {
	Addr         string
	TxPowerLevel int
	Rssi         int
	RawMfrData   string
}

// stringer for observations
func (o FBeacon) String() string {
	// metadata, Temp:value, Acc:value
	return fmt.Sprintf("(%v):[%v](%v), Temp: %v, Acc: (%v, %v, %v)", o.Timestamp.Format(time.Stamp), o.BtData.Addr, o.BtData.Rssi, o.Measurements.Temp, o.Measurements.XAcc, o.Measurements.YAcc, o.Measurements.ZAcc)
}
