package hub

// Includes logic to activate BLE and snatch fujitsu tag readings from the air.
// Runs as it's own routine from hub.go and uses an exported channel to send readings to analysis routine.

// https://towardsdatascience.com/spelunking-bluetooth-le-with-go-c2cff65a7aca
// https://www.thepolyglotdeveloper.com/2018/02/scan-ble-ibeacon-devices-golang-raspberry-pi-zero-w/
// https://github.com/muka/go-bluetooth

// TODO Badly needs refactored, this is a poc using copied code as I learn the language.
// The fuji math is probably wrong

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"regexp"
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
		d.Scan([]gatt.UUID{}, true)
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
		obs.temp = calcTemp(fujiHexToFloat(pkt[1]))
		obs.xAcc = calcAcc(fujiHexToFloat(pkt[2]))
		obs.yAcc = calcAcc(fujiHexToFloat(pkt[3]))
		obs.zAcc = calcAcc(fujiHexToFloat(pkt[4]))
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
// Use binary package to create Uint16 from flipped bits
// Cast to float32 and return
func fujiHexToFloat(hval string) float32 {
	flipped := make([]byte, 4)
	b, _ := hex.DecodeString(hval)
	fmt.Printf("Input: %v, Decoded: %v\n\n", hval, b)
	flipped[0] = b[2]
	flipped[1] = b[3]
	flipped[2] = b[0]
	flipped[3] = b[1]
	return float32(binary.BigEndian.Uint16(flipped))
}

func calcTemp(raw float32) float32 {
	return (((raw / 333.87) + 21.0) * 9.0 / 5.0) + 32
}

func calcAcc(raw float32) float32 {
	return raw / 2048.0
}

//obs.xAcc = float32(binary.BigEndian.Uint16([]byte(pkt[2])) / 2048.0)
//obs.yAcc = float32(binary.BigEndian.Uint16([]byte(pkt[3])) / 2048.0)
//obs.zAcc = float32(binary.BigEndian.Uint16([]byte(pkt[4])) / 2048.0)
