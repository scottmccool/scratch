package hub

// https://towardsdatascience.com/spelunking-bluetooth-le-with-go-c2cff65a7aca
// https://www.thepolyglotdeveloper.com/2018/02/scan-ble-ibeacon-devices-golang-raspberry-pi-zero-w/
// https://github.com/muka/go-bluetooth
// TODO Badly needs refactored, this is a poc using copied code as I learn the language.

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

// ScanFuji - Scan BTLE for a Fujitsu beacon, return it as a payload
func ScanFuji() string {
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
		d.Scan([]gatt.UUID{}, false)
		return
	default:
		d.StopScanning()
	}
}

func onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	reading, err := makeFujiTag(p, a, rssi)
	if err != nil {
		fmt.Println("ID: ", p.ID(), "Error: ", err)
		//		fmt.Printf("ID: %s __not a fuji tag__\n", p.ID())
	} else {
		fmt.Println("***  Found Fuji Tag ***")
		fmt.Printf("%+v \n", reading)
		fmt.Println("***************")
		//		rawc <- reading
	}
}

// FBeacon - a Fujitsu tag
type FBeacon struct {
	hexString       string
	tagId           string
	rawData         []byte
	readingsOnlyHex string
	//uuid  string
	//major uint16
	//minor uint16
}

// makeFujiTag
// A fujitsu tag packet as read by gatt looks like this:
//p.ID(): E2:94:B4:AF:93:13
//p.Name():
//a.LocalName:
// a.TxPowerLevel: 0
// hex.EncodeToString(a.ManufacturerData): 590001000300030034065bfbb1febb06
// a.ServiceData: []
func makeFujiTag(p gatt.Peripheral, a *gatt.Advertisement, rssi int) (FBeacon, error) {
	var reading FBeacon
	vendorString, _ := hex.DecodeString("5900010003000300") // MfrData starts with this for our tags
	if len(vendorString) <= len(a.ManufacturerData) {
		if bytes.Equal(vendorString, a.ManufacturerData[0:len(vendorString)]) {
			reading.rawData = a.ManufacturerData
			reading.hexString = hex.EncodeToString(a.ManufacturerData)
			reading.tagId = p.ID()
			reading.readingsOnlyHex = hex.EncodeToString(a.ManufacturerData[len(a.ManufacturerData)-len(vendorString):])
			return reading, nil
		}
	}
	return reading, errors.New("Not a fuji tag")
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
