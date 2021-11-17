package devicevalue

import (
	"time"
)

// ###########################################################################
// ###########################################################################
// ###########################################################################
// ###########################################################################
// ###########################################################################

// DeviceValue ...
type DeviceValue struct {
	Value     int       `json:"raw"`
	Formatted string    `json:"formatted"`
	Stamp     time.Time `json:"stamp"`
}

// ---------------------------------------------------------------------------

// IDeviceValue is a fake device
type IDeviceValue interface {
	GetValue() int
	GetFormatted() string
	GetStamp() time.Time
}

// ###########################################################################

// Init ...
func Init(v *DeviceValue, value int, formatted string, stamp time.Time) {
	v.Value = value
	v.Formatted = formatted
	v.Stamp = stamp
}

// InitFromDeviceValue ...
func InitFromDeviceValue(v *DeviceValue, src IDeviceValue) {
	v.Value = src.GetValue()
	v.Formatted = src.GetFormatted()
	v.Stamp = src.GetStamp()
}

// ###########################################################################

// GetValue ...
func (v DeviceValue) GetValue() int { return v.Value }

// GetFormatted ...
func (v DeviceValue) GetFormatted() string { return v.Formatted }

// GetStamp ...
func (v DeviceValue) GetStamp() time.Time { return v.Stamp }
