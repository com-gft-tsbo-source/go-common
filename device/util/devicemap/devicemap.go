package devicemap

import (
	"fmt"

	"github.com/com-gft-tsbo-source/go-common/device/implementation/devicedescriptor"
	"github.com/com-gft-tsbo-source/go-common/device/util/devicenode"
)

// ###########################################################################
// ###########################################################################
// ###########################################################################
// ###########################################################################
// ###########################################################################

// DeviceMap is the abstraction of the device map
type DeviceMap struct {
	ByDeviceAddress map[string]*devicenode.DeviceNode
	sequence        map[int]*devicenode.DeviceNode
	maxSequence     int
}

// ---------------------------------------------------------------------------

// IDeviceMap implements the device map
type IDeviceMap interface {
	Add(dev *devicedescriptor.DeviceDescriptor, urlDevice string, urlMeasure string, urlStatus string) (*devicenode.DeviceNode, error)
	Get(DeviceAddress string) (*devicenode.DeviceNode, error)
	Delete(DeviceAddress string) (*devicenode.DeviceNode, error)
	Len() int
}

// ###########################################################################

// InitDeviceMap Create a new DeviceMap
func InitDeviceMap(d *DeviceMap) {
	d.ByDeviceAddress = map[string]*devicenode.DeviceNode{}
	d.sequence = map[int]*devicenode.DeviceNode{}
	d.maxSequence = 0
}

// ###########################################################################

// Add adds a new entry to the DeviceMap
func (dm *DeviceMap) Add(node *devicenode.DeviceNode) (*devicenode.DeviceNode, error) {
	foundNode, found := dm.ByDeviceAddress[node.DeviceAddress]
	if found {
		return foundNode, &OpErrorDuplicate{fmt.Sprintf("Duplicate device '%s'", node.DeviceAddress)}
	}
	dm.maxSequence = dm.maxSequence + 1
	sequence := dm.maxSequence
	node.SetSequence(sequence)
	dm.ByDeviceAddress[node.DeviceAddress] = node
	dm.sequence[sequence] = node
	return node, nil
}

// ---------------------------------------------------------------------------

// Get returns a device by its name
func (dm *DeviceMap) Get(DeviceAddress string) (*devicenode.DeviceNode, error) {
	foundNode, found := dm.ByDeviceAddress[DeviceAddress]
	if !found {
		return nil, &OpErrorNotFound{fmt.Sprintf("Device '%s' not found", DeviceAddress)}
	}
	return foundNode, nil
}

// ---------------------------------------------------------------------------

// Delete returns a device by its name
func (dm *DeviceMap) Delete(DeviceAddress string) (*devicenode.DeviceNode, error) {
	foundNode, found := dm.ByDeviceAddress[DeviceAddress]
	if !found {
		return nil, &OpErrorNotFound{fmt.Sprintf("Device '%s' not found", DeviceAddress)}
	}
	delete(dm.ByDeviceAddress, DeviceAddress)
	delete(dm.sequence, foundNode.GetSequence())
	return foundNode, nil
}

// ---------------------------------------------------------------------------

func (dm *DeviceMap) len() int {
	return len(dm.ByDeviceAddress)
}
