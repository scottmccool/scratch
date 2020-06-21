package hub

// Includes logic to activate BLE and snatch fujitsu tag readings from the air.
// Runs as it's own routine from hub.go and uses an exported channel to send readings to analysis routine.
// Uses paypal's gatt for the btle heavy lifting, and barely customizes anything

// https://towardsdatascience.com/spelunking-bluetooth-le-with-go-c2cff65a7aca
// https://www.thepolyglotdeveloper.com/2018/02/scan-ble-ibeacon-devices-golang-raspberry-pi-zero-w/
// https://github.com/muka/go-bluetooth

// TODO Badly needs refactored, this is a poc using copied code as I learn the language.
// The fuji math is probably wrong in an ugly way.

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/scottmccool/FBeacon/readings"

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
	d.Handle(gatt.PeripheralDiscovered(sniffFujiTag))
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

// sniffFujiTag Handler to gatt discovery, see if the BLE reading is one of interest and return a measurement if possible
func sniffFujiTag(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	hexMfr := hex.EncodeToString(a.ManufacturerData)
	beacon, err := extractFujiTag(hexMfr)
	if err != nil {
		return
	}
	beacon.BtData.Addr = p.ID()
	beacon.BtData.Rssi = rssi
	beacon.BtData.RawMfrData = hexMfr
	beacon.BtData.TxPowerLevel = a.TxPowerLevel
	Rawc <- beacon // Publish to channel for analysis by rest of hub; we are done

}

//We care about adv type 255 manufacturer data, the following regex from python code shows format
//PACKET_DATA_REGEX = re.compile(r'010003000300(?P<temperature>.{4})(?P<x_acc>.{4})(?P<y_acc>.{4})(?P<z_acc>.{4})$')
// Bytes need flipped
// Temp formula from fuji: (((_unpack_value(_flip_bytes(hex_temperature)) / 333.87) + 21.0) * 9.0 / 5.0) + 32
//   So: hex_temp=c401 raw, flip the bytes to 01c4, cast it to an int16 (return 2^^16-val if exceeds), then do the math above results in t_f=72.2368
// Acc formula from fuji: _unpack_value(_flip_bytes(hex_accel)) / 2048.0
// makeFujiTag() Builds a FBeacon object with decoded data from a tag.  Writes it to channel.
func extractFujiTag(hexMfr string) (beacon readings.FBeacon, err error) {
	re := regexp.MustCompile(`010003000300(?P<temperature>.{4})(?P<x_acc>.{4})(?P<y_acc>.{4})(?P<z_acc>.{4})$`)
	pkt := re.FindStringSubmatch(hexMfr)
	if len(pkt) == 5 { // 0 always empty
		beacon.Timestamp = time.Now()
		beacon.Measurements.Temp = calcTemp(fujiHexToUInt(pkt[1]))
		beacon.Measurements.XAcc = calcAcc(fujiHexToUInt(pkt[2]))
		beacon.Measurements.YAcc = calcAcc(fujiHexToUInt(pkt[3]))
		beacon.Measurements.ZAcc = calcAcc(fujiHexToUInt(pkt[4]))
		//j, _ := json.Marshal(beacon)
		//fmt.Printf("--Input: %v\n--Output: %v\n", hexMfr, string(j))
	} else {
		err = errors.New("Does not contain Fujitsu beacon data")
	}
	return beacon, err
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
	//fmt.Printf("--Hex flipper, string: %v flipped to int: %v\n", hval, rv)
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
