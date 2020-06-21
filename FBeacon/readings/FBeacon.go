package readings

import (
	"fmt"
	"time"
)

// Our sensor data reading object(s)

// FBeacon a fujitsu tag measuring temp, x_acc, y_acc, z_acc along with bluetooth metadata
type FBeacon struct {
	timestamp    time.Time
	bt_data      Bt_data
	measurements Measurements
}

type measurements struct {
	temp float32
	xAcc float32
	yAcc float32
	zAcc float32
}

type bt_data struct {
	addr         string
	txPowerLevel int
	rssi         int
	rawMfrData   string
}

// stringer for observations
func (o FBeacon) String() string {
	// metadata, Temp:value, Acc:value
	return fmt.Sprintf("(%v):[%v](%v), Temp: %v, Acc: (%v, %v, %v)", o.timestamp.Format(time.Stamp), o.bt_data.addr, o.bt_data.rssi, o.measurements.temp, o.measurements.xAcc, o.measurements.yAcc, o.measurements.zAcc)
}
