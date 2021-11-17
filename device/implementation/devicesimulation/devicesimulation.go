package devicesimulation

import (
	"math/rand"
	"time"
)

// ###########################################################################
// ###########################################################################
// ###########################################################################
// ###########################################################################
// ###########################################################################

// DeviceSimulation ...
type DeviceSimulation struct {
	SeededRand  *rand.Rand `json:"-"`
	SimValue    int        `json:"-"`
	SimMedian   int        `json:"simMedian"`
	SimVariance int        `json:"simVariance"`
	SimInterval int        `json:"simInterval"`
	SimAmount   int        `json:"simAmount"`
}

// ---------------------------------------------------------------------------

// IDeviceSimulation ...
type IDeviceSimulation interface {
	GetValue() int
	GetMedian() int
	GetVariance() int
	GetInterval() int
	GetAmount() int
	Simulate()
	TranslateValue(i int)
}

// ###########################################################################
// New ...

// Init ...
func Init(v *DeviceSimulation, seededRand *rand.Rand, median int, variance int, value int, interval int, amount int) {
	if seededRand == nil {
		seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	v.SeededRand = seededRand
	v.SimMedian = median
	v.SimVariance = variance
	v.SimValue = value
	v.SimInterval = interval
	v.SimAmount = amount
}

// InitFromDeviceSimulation ...
func InitFromDeviceSimulation(v *DeviceSimulation, src IDeviceSimulation) {
	v.SeededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	v.SimMedian = src.GetMedian()
	v.SimVariance = src.GetVariance()
	v.SimValue = src.GetValue()
	v.SimInterval = src.GetInterval()
	v.SimAmount = src.GetAmount()
}

// ###########################################################################

// Simulate slightly alters the sim value
func (v *DeviceSimulation) Simulate() {
	v.SimValue = v.SimValue + v.SeededRand.Intn(100) - 50
	if v.SimValue > v.SimMedian+v.SimVariance {
		v.SimValue = v.SimMedian + v.SimVariance
	} else if v.SimValue < v.SimMedian-v.SimVariance {
		v.SimValue = v.SimMedian - v.SimVariance
	}
}

// TranslateValue Set the value from a given integeger
func (v *DeviceSimulation) TranslateValue(i int) {

	var amount int

	if i > 100 {
		i = 100
	} else if i < 0 {
		i = 0
	}
	i = i - 50

	amount = v.SimAmount * i / 100

	if amount == 0 {
		amount = 1
	}

	v.SimValue = v.SimValue + amount

	if v.SimValue > v.SimMedian+v.SimVariance {
		v.SimValue = v.SimMedian + v.SimVariance
	} else if v.SimValue < v.SimMedian-v.SimVariance {
		v.SimValue = v.SimMedian - v.SimVariance
	}
}

// ---------------------------------------------------------------------------

// GetValue returns the current value
func (v *DeviceSimulation) GetValue() int { return v.SimValue }

// ---------------------------------------------------------------------------

// GetMedian ...
func (v *DeviceSimulation) GetMedian() int { return v.SimMedian }

// ---------------------------------------------------------------------------

// GetVariance ...
func (v *DeviceSimulation) GetVariance() int { return v.SimVariance }

// ---------------------------------------------------------------------------

// GetVariance ...
func (v *DeviceSimulation) GetAmount() int { return v.SimAmount }

// ---------------------------------------------------------------------------

// GetInterval ...
func (v *DeviceSimulation) GetInterval() int { return v.SimInterval }
