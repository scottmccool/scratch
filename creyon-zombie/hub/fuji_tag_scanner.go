package hub

// Includes logic to activate BLE and snatch fujitsu tag readings from the air.
// Runs as it's own routine from hub.go and uses an exported channel to send readings to analysis routine.

// https://towardsdatascience.com/spelunking-bluetooth-le-with-go-c2cff65a7aca
// https://www.thepolyglotdeveloper.com/2018/02/scan-ble-ibeacon-devices-golang-raspberry-pi-zero-w/
// https://github.com/muka/go-bluetooth

// TODO Badly needs refactored, this is a poc using copied code as I learn the language.
// The fuji math is probably wrong in an ugly way.

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

// ScanFuji - Scan BTLE for a Fujitsu beacon, return it as a payload.  Entrypoint.
func ScanFuji() {
	//	return "Mock-reading"
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
	}

	// Register handlers.
	d.Handle(gatt.PeripheralDiscovered(makeFujiTag))
	d.Init(onStateChanged)
	select {}
}

func onStateChanged(d gatt.Device, s gatt.State) {
	fmt.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("scanning...")
		d.Scan([]gatt.UUID{}, true) // Ignore second ad from a given device; assumes we are being externally activated to snatch at given resolution.
		return
	default:
		d.StopScanning()
	}
}

//We care about adv type 255 manufacturer data, the following regex from python code shows format
//PACKET_DATA_REGEX = re.compile(r'010003000300(?P<temperature>.{4})(?P<x_acc>.{4})(?P<y_acc>.{4})(?P<z_acc>.{4})$')
// Bytes need flipped
// Temp formula from fuji: (((_unpack_value(_flip_bytes(hex_temperature)) / 333.87) + 21.0) * 9.0 / 5.0) + 32
//   So: hex_temp=c401 raw, flip the bytes to 01c4, cast it to an int16 (return 2^^16-val if exceeds), then do the math above results in t_f=72.2368
// Acc formula from fuji: _unpack_value(_flip_bytes(hex_accel)) / 2048.0
//def _flip_bytes(hex_bytes):
//  return ''.join(map(lambda pr: ''.join(pr), each_slice(2, list(hex_bytes)))[::-1])
// makeFujiTag() Builds a FBeacon object with decoded data from a tag.  Writes it to channel.
func makeFujiTag(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	var obs FBeacon
	re := regexp.MustCompile(`010003000300(?P<temperature>.{4})(?P<x_acc>.{4})(?P<y_acc>.{4})(?P<z_acc>.{4})$`)
	hexMfr := hex.EncodeToString(a.ManufacturerData)
	pkt := re.FindStringSubmatch(hexMfr)
	if len(pkt) == 5 { // 0 always empty
		obs.temp = calcTemp(fujiHexToUInt(pkt[1]))
		obs.xAcc = calcAcc(fujiHexToUInt(pkt[2]))
		obs.yAcc = calcAcc(fujiHexToUInt(pkt[3]))
		obs.zAcc = calcAcc(fujiHexToUInt(pkt[4]))
		obs.addr = p.ID()
		obs.txPowerLevel = a.TxPowerLevel
		obs.rawMfrData = hexMfr
		obs.timestamp = time.Now()
		obs.rssi = rssi
		Rawc <- obs // Publish to channel for analysis by rest of hub logic, we're done
	}
}

// Turn a 4 byte hex string from the advertisement to a float32
// Flip bits so that c401 becomes 01c4
// Use binary package to create Uint16 from flipped bits (//TODO verify this and encoding)
// Cast to float32 and return for inclusion in an FBeacon
func fujiHexToUInt(hval string) uint16 {
	//orig_bytes, _ := hex.DecodeString(hval)
	chars := strings.Split(hval, "")
	flipped := make([]string, 4)
	flipped[0] = chars[2]
	flipped[1] = chars[3]
	flipped[2] = chars[0]
	flipped[3] = chars[1]
	flippedS := strings.Join(flipped, "")
	hFlipped, _ := hex.DecodeString(flippedS)
	rv := binary.BigEndian.Uint16(hFlipped)
	//fmt.Printf("RV: %v from Flipped %v(%v) to %v(%v)\n", rv, hval, orig_bytes, flipped_s, h_flipped)
	return rv
}

// calcTemp convert a decoeded raw vlue into temperature(f) using Fuji-provided math
func calcTemp(raw uint16) float32 {
	f := float32(raw)
	return (((f / 333.87) + 21.0) * 9.0 / 5.0) + 32
}

// calcAcc convert a decoded raw value into an accelerometer reading using Fuji-provided math
func calcAcc(raw uint16) float32 {
	f := float32(raw)
	return f / 2048.0
}
