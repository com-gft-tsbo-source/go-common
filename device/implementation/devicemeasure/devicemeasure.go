package devicemeasure

import (
	"time"

	"github.com/com-gft-tsbo-source/go-common/device/implementation/devicedescriptor"
	"github.com/com-gft-tsbo-source/go-common/device/implementation/devicevalue"
)

// ###########################################################################
// ###########################################################################
// ###########################################################################
// ###########################################################################
// ###########################################################################

// DeviceMeasure Encapsulates the reploy of ms-thermometer
type DeviceMeasure struct {
	devicedescriptor.DeviceDescriptor
	devicevalue.DeviceValue
}

// ---------------------------------------------------------------------------

// IDeviceMeasure ...
type IDeviceMeasure interface {
	devicedescriptor.IDeviceDescriptor
	devicevalue.IDeviceValue
}

// ###########################################################################

// Init ...
func Init(d *DeviceMeasure, deviceType string, deviceAddress string, unit string,
	value int, formatted string, stamp time.Time, interval int) {
	devicedescriptor.Init(&d.DeviceDescriptor, deviceType, deviceAddress, unit)
	devicevalue.Init(&d.DeviceValue, -1, "", time.Now())
}

// InitFromDeviceMeasure ...
func InitFromDeviceMeasure(d *DeviceMeasure, src IDeviceMeasure) {
	devicedescriptor.InitFromDeviceDescriptor(&d.DeviceDescriptor, src)
	devicevalue.Init(&d.DeviceValue, -1, "", time.Now())
}
