package hub

import (
	"testing"
)

func TestextractFujiTag(t *testing.T) {

}

func TestFujiHexToUInt(t *testing.T) {
	cases := []struct {
		in   string
		want uint16
	}{
		{"5401", 340},
		{"2c00", 44},
		{"1300", 19},
		{"8f08", 2191},
	}
	{
		for _, c := range cases {
			got := fujiHexToUInt(c.in)
			if got != c.want {
				t.Errorf("fujiHexToUnit(%v) == %v, want %v", c.in, got, c.want)
			}
		}
	}
}

func TestCalcTemp(t *testing.T) {
	cases := []struct {
		in   string
		want uint16
	}{
		{"5401", 340},
		{"2c00", 44},
		{"1300", 19},
		{"8f08", 2191},
	}
	{
		for _, c := range cases {
			got := fujiHexToUInt(c.in)
			if got != c.want {
				t.Errorf("fujiHexToUnit(%v) == %v, want %v", c.in, got, c.want)
			}
		}
	}
}

func TestCalcAcc(t *testing.T) {
	ans := calcAcc(uint16(140))
	if ans != 70.2 {
		t.Errorf("%v converted to %v, expected %v", uint(140), ans, float32(12.0))
	}
}

//Hex flipper, string: 5401 flipped to int: 340
//Hex flipper, string: 2c00 flipped to int: 44
//Hex flipper, string: 1300 flipped to int: 19
//Hex flipper, string: 8f08 flipped to int: 2191
//Input: 590001000300030054012c0013008f08
//Output: {"Timestamp":"2020-06-20T21:07:05.657304784-07:00","BtData":{"Addr":"","TxPowerLevel":0,"Rssi":0,"RawMfrData":""},"Measurements":{"Temp":71.63305,"XAcc":0.021484375,"YAcc":0.009277344,"ZAcc":1.0698242}}
//Hex flipper, string: 1501 flipped to int: 277
//Hex flipper, string: 5d00 flipped to int: 93
//Hex flipper, string: 1800 flipped to int: 24
//Hex flipper, string: 4308 flipped to int: 2115
//Input: 590001000300030015015d0018004308
//Output: {"Timestamp":"2020-06-20T21:07:05.951338914-07:00","BtData":{"Addr":"","TxPowerLevel":0,"Rssi":0,"RawMfrData":""},"Measurements":{"Temp":71.293396,"XAcc":0.045410156,"YAcc":0.01171875,"ZAcc":1.0327148}}
//Hex flipper, string: 5401 flipped to int: 340
//Hex flipper, string: 3500 flipped to int: 53
//Hex flipper, string: 0c00 flipped to int: 12
//Hex flipper, string: 8f08 flipped to int: 2191
//Input: 5900010003000300540135000c008f08
//Output: {"Timestamp":"2020-06-20T21:07:06.662349876-07:00","BtData":{"Addr":"","TxPowerLevel":0,"Rssi":0,"RawMfrData":""},"Measurements":{"Temp":71.63305,"XAcc":0.025878906,"YAcc":0.005859375,"ZAcc":1.0698242}}
//Hex flipper, string: 0501 flipped to int: 261
//Hex flipper, string: 5400 flipped to int: 84
//Hex flipper, string: 1a00 flipped to int: 26
//Hex flipper, string: 4b08 flipped to int: 2123
//Input: 5900010003000300050154001a004b08
//Output: {"Timestamp":"2020-06-20T21:07:06.953335342-07:00","BtData":{"Addr":"","TxPowerLevel":0,"Rssi":0,"RawMfrData":""},"Measurements":{"Temp":71.20714,"XAcc":0.041015625,"YAcc":0.0126953125,"ZAcc":1.0366211}}
//Hex flipper, string: 8401 flipped to int: 388
//Hex flipper, string: 3100 flipped to int: 49
//Hex flipper, string: 1300 flipped to int: 19
//Hex flipper, string: 8e08 flipped to int: 2190
//Input: 59000100030003008401310013008e08
//Output: {"Timestamp":"2020-06-20T21:07:07.665320827-07:00","BtData":{"Addr":"","TxPowerLevel":0,"Rssi":0,"RawMfrData":""},"Measurements":{"Temp":71.89183,"XAcc":0.023925781,"YAcc":0.009277344,"ZAcc":1.0693359}}
//Hex flipper, string: d500 flipped to int: 213
//Hex flipper, string: 5800 flipped to int: 88
//Hex flipper, string: 1900 flipped to int: 25
//Hex flipper, string: 4408 flipped to int: 2116
//Input: 5900010003000300d500580019004408
//Output: {"Timestamp":"2020-06-20T21:07:07.961355475-07:00","BtData":{"Addr":"","TxPowerLevel":0,"Rssi":0,"RawMfrData":""},"Measurements":{"Temp":70.94835,"XAcc":0.04296875,"YAcc":0.012207031,"ZAcc":1.0332031}}
//
