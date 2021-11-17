package device

import (
	"math/rand"
	"time"

	"github.com/com-gft-tsbo-source/go-common/device/implementation/devicedescriptor"
	"github.com/com-gft-tsbo-source/go-common/device/implementation/devicesimulation"
	"github.com/com-gft-tsbo-source/go-common/device/implementation/devicevalue"
)

// ###########################################################################
// ###########################################################################
// ###########################################################################
// ###########################################################################
// ###########################################################################

// Device keeps all info to simulate a fake device
type Device struct {
	devicedescriptor.DeviceDescriptor
	devicesimulation.DeviceSimulation
}

// ---------------------------------------------------------------------------

// IDevice ...
type IDevice interface {
	devicedescriptor.IDeviceDescriptor
	devicesimulation.IDeviceSimulation
	FillDeviceValue(*devicevalue.DeviceValue)
}

// ###########################################################################

// InitDevice ...
func InitDevice(d *Device, deviceType string, deviceAddress string, unit string,
	seededRand *rand.Rand, median int, variance int, interval int, amount int) {
	devicedescriptor.Init(&d.DeviceDescriptor, deviceType, deviceAddress, unit)
	devicesimulation.Init(&d.DeviceSimulation, seededRand, median, variance, -1, interval, amount)
}

// ---------------------------------------------------------------------------

// InitThermometer ...
func InitThermometer(d *Device, deviceAddress string, median int, variance int, interval int, amount int) {

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	if median == -1 {
		median = seededRand.Intn(11)*100 + 1500
	}

	if variance == -1 {
		variance = seededRand.Intn(3)*100 + 100
	}

	if interval == -1 {
		interval = seededRand.Intn(2000) + 250
	}

	if amount == -1 {
		amount = seededRand.Intn(200)
	}

	InitDevice(d, "thermometer", deviceAddress, "C", seededRand, median, variance, interval, amount)
	d.Simulate()
}

// ---------------------------------------------------------------------------

// InitHygrometer ...
func InitHygrometer(d *Device, deviceAddress string, median int, variance int, interval int, amount int) {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	if median == -1 {
		median = seededRand.Intn(60)*100 + 2000
	}

	if variance == -1 {
		variance = seededRand.Intn(20)*100 + 100
	}

	if interval == -1 {
		interval = seededRand.Intn(4)*1000 + 2000
	}

	if amount == -1 {
		amount = seededRand.Intn(200)
	}

	InitDevice(d, "hygrometer", deviceAddress, "%", seededRand, median, variance, interval, amount)
	d.Simulate()
}

// ###########################################################################

// ---------------------------------------------------------------------------

// FillDeviceValue ...
func (d *Device) FillDeviceValue(v *devicevalue.DeviceValue) {
	value := d.GetValue()
	fmtValue := d.FormatValue(value)
	devicevalue.Init(v, value, fmtValue, time.Now())
}
