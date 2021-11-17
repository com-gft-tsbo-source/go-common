package devicedescriptor

import (
	"fmt"
)

// ###########################################################################
// ###########################################################################
// ###########################################################################
// ###########################################################################
// ###########################################################################

// DeviceDescriptor is the minimal description of a device
type DeviceDescriptor struct {
	DeviceType    string `json:"type"`
	DeviceAddress string `json:"address"`
	Unit          string `json:"unit"`
}

// ---------------------------------------------------------------------------

// IDeviceDescriptor is a fake device
type IDeviceDescriptor interface {
	GetDeviceType() string
	GetDeviceAddress() string
	GetUnit() string
	FormatValue(value int) string
}

// ###########################################################################

// Init ...
func Init(d *DeviceDescriptor, DeviceType string, DeviceAddress string, unit string) {
	d.DeviceType = DeviceType
	d.DeviceAddress = DeviceAddress
	d.Unit = unit
}

// InitFromDeviceDescriptor ...
func InitFromDeviceDescriptor(d *DeviceDescriptor, src IDeviceDescriptor) {
	d.DeviceType = src.GetDeviceType()
	d.DeviceAddress = src.GetDeviceAddress()
	d.Unit = src.GetUnit()
}

// ###########################################################################

// GetDeviceType ...
func (d DeviceDescriptor) GetDeviceType() string { return d.DeviceType }

// ---------------------------------------------------------------------------

// GetDeviceAddress ...
func (d DeviceDescriptor) GetDeviceAddress() string { return d.DeviceAddress }

// ---------------------------------------------------------------------------

// GetUnit ...
func (d DeviceDescriptor) GetUnit() string { return d.Unit }

// ---------------------------------------------------------------------------

// FormatValue ...
func (d DeviceDescriptor) FormatValue(value int) string {
	return fmt.Sprintf("%4.2f%s", float64(value)/100, d.Unit)
}
