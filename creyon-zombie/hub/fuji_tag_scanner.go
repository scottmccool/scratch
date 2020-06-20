package hub

// https://towardsdatascience.com/spelunking-bluetooth-le-with-go-c2cff65a7aca
// https://www.thepolyglotdeveloper.com/2018/02/scan-ble-ibeacon-devices-golang-raspberry-pi-zero-w/
// https://github.com/muka/go-bluetooth
// TODO Badly needs refactored, this is a poc using copied code as I learn the language.

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

// ScanFuji - Scan BTLE for a Fujitsu beacon, return it as a payload
func ScanFuji() {
	//	return "Mock-reading"
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
	}

	// Register handlers.
	d.Handle(gatt.PeripheralDiscovered(onPeripheralDiscovered))
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

func onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	reading, err := makeFujiTag(p, a, rssi)
	if err != nil {
		//		fmt.Println("ID: ", p.ID(), "Msg: ", err)
		//		fmt.Printf("ID: %s __not a fuji tag__\n", p.ID())
	} else {
		Rawc <- reading
		fmt.Printf("Grabbed fuji packet from %v, sleeping 1 second\n", p.ID())
	}
}

// makeFujiTag
// A fujitsu tag packet as read by gatt looks like this:
//p.ID(): E2:94:B4:AF:93:13
//p.Name():
//a.LocalName:
// hex.EncodeToString(a.ManufacturerData): 590001000300030034065bfbb1febb06

//We care about adv type 255 manufacturer data, the following regex from python code shows format
//PACKET_DATA_REGEX = re.compile(r'010003000300(?P<temperature>.{4})(?P<x_acc>.{4})(?P<y_acc>.{4})(?P<z_acc>.{4})$')
// Bytes need flipped
// Temp formula from fuji: (((_unpack_value(_flip_bytes(hex_temperature)) / 333.87) + 21.0) * 9.0 / 5.0) + 32
// Acc formula from fuji: _unpack_value(_flip_bytes(hex_accel)) / 2048.0
//def _flip_bytes(hex_bytes):
//  return ''.join(map(lambda pr: ''.join(pr), each_slice(2, list(hex_bytes)))[::-1])
func makeFujiTag(p gatt.Peripheral, a *gatt.Advertisement, rssi int) (FBeacon, error) {
	var obs FBeacon
	re := regexp.MustCompile(`010003000300(?P<temperature>.{4})(?P<x_acc>.{4})(?P<y_acc>.{4})(?P<z_acc>.{4})$`)
	hexMfr := hex.EncodeToString(a.ManufacturerData)
	pkt := re.FindStringSubmatch(hexMfr)
	if len(pkt) == 5 { // 0 always empty
		obs.temp = float32((((float32(binary.LittleEndian.Uint16([]byte(pkt[1]))) / 333.87) + 21.0) * 9.0 / 5.0) + 32)
		obs.xAcc = float32(binary.LittleEndian.Uint16([]byte(pkt[2])) / 2048.0)
		obs.yAcc = float32(binary.LittleEndian.Uint16([]byte(pkt[2])) / 2048.0)
		obs.zAcc = float32(binary.LittleEndian.Uint16([]byte(pkt[2])) / 2048.0)
		obs.addr = p.ID()
		obs.txPowerLevel = a.TxPowerLevel
		obs.rawMfrData = a.ManufacturerData
		obs.timestamp = time.Now()
		return obs, nil
	}
	return obs, errors.New("Not a fujitsu beacon")
}

// Used for reverse engineering ble ads
func printAdData(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	var flipped []byte
	vendorStringBytes := []byte("10003000300")
	fmt.Printf("\nPeripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())
	fmt.Printf("\n-----\n%+v\n----\n", a)
	fmt.Println("  Local Name        =", a.LocalName)
	fmt.Println("  TX Power Level    =", a.TxPowerLevel)
	fmt.Println("  Manufacturer Data =", hex.EncodeToString(a.ManufacturerData))
	if len(a.ManufacturerData) >= len(vendorStringBytes) {
		for _, bit := range a.ManufacturerData {
			flipped = append(flipped, ^bit)
		}
		fmt.Println("** May be a fuji packet **")
		fmt.Println("              Raw: ", a.ManufacturerData)
		fmt.Println("          Flipped: ", flipped)
		fmt.Println("              Hex: ", hex.EncodeToString(a.ManufacturerData))
		fmt.Println(bytes.Split(a.ManufacturerData, []byte(":")))
		fmt.Println(bytes.Split(flipped, []byte(":")))
		fmt.Println("          Flipped: ", hex.EncodeToString(flipped))
		fmt.Println("** End Fuji details **")
	}
	fmt.Println("  Service Data      =", a.ServiceData)
}
