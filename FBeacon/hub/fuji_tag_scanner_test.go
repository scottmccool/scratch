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
		{"f500", 245},
		{"4401", 324},
		{"2c00", 44},
		{"5e00", 94},
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
		in   uint16
		want float32
	}{
		{245, 71.12087},
		{70, 70.17739},
	}
	{
		for _, c := range cases {
			got := calcTemp(c.in)
			if got != c.want {
				t.Errorf("calcTemp(%v) == %v, want %v", c.in, got, c.want)
			}
		}
	}
}

func TestCalcAcc(t *testing.T) {
	cases := []struct {
		in   uint16
		want float32
	}{
		{44, 0.021484375},
		{94, 0.045898438},
	}
	{
		for _, c := range cases {
			got := calcAcc(c.in)
			if got != c.want {
				t.Errorf("calcTemp(%v) == %v, want %v", c.in, got, c.want)
			}
		}
	}
}

/*
Test data:
Warning, generated from implemention (;
--                                                                                                                          vendor    010003000300
{"Timestamp":"2020-06-20T21:25:49.1262864-07:00","BtData":{"Addr":"E2:94:B4:AF:93:13","TxPowerLevel":0,"Rssi":-73,  "RawMfrData":"590001000300030044012c0011007b08"},"Measurements":{"Temp":71.54679,"XAcc":0.021484375,"YAcc":0.008300781,"ZAcc":1.0600586}}
{"Timestamp":"2020-06-20T21:25:49.344215004-07:00","BtData":{"Addr":"D6:4F:DE:CF:63:99","TxPowerLevel":0,"Rssi":-77,"RawMfrData":"5900010003000300f5005e0018005208"},"Measurements":{"Temp":71.12087,"XAcc":0.045898438,"YAcc":0.01171875,"ZAcc":1.0400391}}
{"Timestamp":"2020-06-20T21:25:50.132079399-07:00","BtData":{"Addr":"E2:94:B4:AF:93:13","TxPowerLevel":0,"Rssi":-70,"RawMfrData":"5900010003000300440132000f009a08"},"Measurements":{"Temp":71.54679,"XAcc":0.024414062,"YAcc":0.0073242188,"ZAcc":1.0751953}}
{"Timestamp":"2020-06-20T21:25:51.135344492-07:00","BtData":{"Addr":"E2:94:B4:AF:93:13","TxPowerLevel":0,"Rssi":-74,"RawMfrData":"5900010003000300140121000d008608"},"Measurements":{"Temp":71.28801,"XAcc":0.016113281,"YAcc":0.0063476562,"ZAcc":1.0654297}}

--
*/
