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
<<<<<<< HEAD
	return fmt.Sprintf("(%v):[%v](%v), Temp: %v, Acc: (%v, %v, %v)", o.Timestamp.Format(time.Stamp), o.BtData.Addr, o.BtData.Rssi, o.Measurements.Temp, o.Measurements.XAcc, o.Measurements.YAcc, o.Measurements.ZAcc)
=======
	return fmt.Sprintf("(%v):[%v](%v), Temp: %v, Acc: (%v, %v, %v)", o.timestamp.Format(time.Stamp), o.bt_data.addr, o.bt_data.rssi, o.measurements.temp, o.measurements.xAcc, o.measurements.yAcc, o.measurements.zAcc)
>>>>>>> 9b47655055a47faf961c28b60604a568da9330c8
}
