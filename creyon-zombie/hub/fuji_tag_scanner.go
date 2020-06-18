package hub

// https://towardsdatascience.com/spelunking-bluetooth-le-with-go-c2cff65a7aca
// https://www.thepolyglotdeveloper.com/2018/02/scan-ble-ibeacon-devices-golang-raspberry-pi-zero-w/
// https://github.com/muka/go-bluetooth
// TODO Badly needs refactored, this is a poc using copied code as I learn the language.

import (
	"encoding/hex"
	"fmt"
	"log"

	//"strconv"

	//	"encoding/hex"

	//	"strings"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

//var tag_readings map[string][string]

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

//        # getScanData returns a tripple from the ScanEntry object
//# the tripple has the advertising type, description and value (adtype, desc, value)
//triples = bleAdvertisement.getScanData()
//# Bluetooth defines AD types https://ianharvey.github.io/bluepy-doc/scanentry.html
//# DREAM only wants adtype = 0xff (0d255) for manufacturer data
//# values is a list where the last element is the manufacturer data
//values = [value for (adtype, desc, value) in triples if adtype == 255]

//strconv.FormatInt(255, 16))
func onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	//x := hex.Dump(a.ManufacturerData)
	x := hex.EncodeToString(a.ManufacturerData)
	//hDataLittle := strconv.FormatUint(binary.LittleEndian.Uint64(a.ManufacturerData))
	//hDataBig := strconv.FormatUint(binary.BigEndian.Uint64(a.ManufacturerData))
	fmt.Printf("\nPeripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())
	fmt.Println("  Local Name        =", a.LocalName)
	fmt.Println("  TX Power Level    =", a.TxPowerLevel)
	fmt.Println("  Manufacturer Data =", a.ManufacturerData)
	fmt.Println("    munged: ", x)
	fmt.Println("  Service Data      =", a.ServiceData)
}

// FBeacon - a Fujitsu tag
type FBeacon struct {
	rawData string
	//uuid  string
	//major uint16
	//minor uint16
}

//Here is the python filter

//PACKET_DATA_REGEX = re.compile(r'010003000300(?P<temperature>.{4})(?P<x_acc>.{4})(?P<y_acc>.{4})(?P<z_acc>.{4})$')

//def decode(packet_hex_string):
//    match = re.search(PACKET_DATA_REGEX, packet_hex_string or '')
//    if not match:
//        return None
//    hex_temperature = match.group('temperature')
//    hex_x_acc = match.group('x_acc')
//    hex_y_acc = match.group('y_acc')
//    hex_z_acc = match.group('z_acc')
//
//    packet = {
//        'temperature': _compute_temperature(hex_temperature),
//        'x_acc': _compute_acceleration(hex_x_acc),
//        'y_acc': _compute_acceleration(hex_y_acc),
//        'z_acc': _compute_acceleration(hex_z_acc)
//    }
//    return packet
//

//def buildPacketFromBle(bleAdvertisement,hci):
//    packet = None
//    try:
//        # get the tag_id (MAC address) and rssi from the BLE advertisement
//        # since the MAC address comes with colons, remove them.
//        tag_id = bleAdvertisement.addr.replace(':', '')
//        rssi = bleAdvertisement.rssi
//
//        # getScanData returns a tripple from the ScanEntry object
//        # the tripple has the advertising type, description and value (adtype, desc, value)
//        triples = bleAdvertisement.getScanData()
//        # Bluetooth defines AD types https://ianharvey.github.io/bluepy-doc/scanentry.html
//        # DREAM only wants adtype = 0xff (0d255) for manufacturer data
//        # values is a list where the last element is the manufacturer data
//        values = [value for (adtype, desc, value) in triples if adtype == 255]
//
//        # TODO: why are we only looking at the last msg recorded?
//        measurements = fujitsu_decoder.decode(values[-1])
//        if measurements is None:
//            return None
//        packet = {
//            "readingType": "bluetooth_fujitsu",
//            "timestamp": int(time.time()),
//            "tag_id": tag_id,
//            "hci": hci,
//            "rssi": rssi,
//            "measurements": measurements
//        }
//    except Exception:
//        # If there's a unicode error, disregard it since it isn't a Fujitsu Tag
//        #print("Discarding on Unicode error (non Fujitsu)")
//        return None
//    return packet
//
