package devicenode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/com-gft-tsbo-source/go-common/device/implementation/devicedescriptor"
	"github.com/com-gft-tsbo-source/go-common/device/implementation/devicevalue"
)

// ###########################################################################
// ###########################################################################
// ###########################################################################
// ###########################################################################
// ###########################################################################

// DeviceNode orders the devices
type DeviceNode struct {
	devicedescriptor.DeviceDescriptor
	devicevalue.DeviceValue
	sequence   int
	isNew      bool
	Cid        string `json:"cid"`
	URLDevice  string `json:"urlDevice"`
	URLMeasure string `json:"urlMeasure"`
	URLStatus  string `json:"urlStatus"`
}

// ---------------------------------------------------------------------------

// IDeviceNode implements the device map
type IDeviceNode interface {
	devicedescriptor.IDeviceDescriptor
	devicevalue.IDeviceValue
	Update() error
	SetSequence(int)
	GetSequence() int
	IsNew() bool
	ClearNew()
	GetURLDevice() string
	GetURLMeasure() string
	GetURLStatus() string
}

// ###########################################################################

// InitDeviceNode creates a new devicenode
func InitDeviceNode(n *DeviceNode, dev devicedescriptor.IDeviceDescriptor, sequence int, cid string, urlDevice string, urlMeasure string, urlStatus string) {
	if dev != nil {
		// devicedescriptor.CopyDeviceDescriptor(&n.DeviceDescriptor, dev)
		devicedescriptor.Init(&n.DeviceDescriptor, dev.GetDeviceType(), dev.GetDeviceAddress(), dev.GetUnit())
	} else {
		devicedescriptor.Init(&n.DeviceDescriptor, "", "", "")
	}
	devicevalue.Init(&n.DeviceValue, -1, "", time.Now())
	n.sequence = sequence
	n.isNew = true
	n.Cid = cid
	n.URLDevice = urlDevice
	n.URLMeasure = urlMeasure
	n.URLStatus = urlStatus
}

// // ---------------------------------------------------------------------------

// FromURL ...
func FromURL(url string) (*DeviceNode, error) {
	httpClient := &http.Client{Timeout: 3 * time.Second}
	r, err := httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	clientID := r.Header.Get("cid")

	if len(clientID) == 0 {
		clientID = r.Header.Get("name")
	}
	if len(clientID) == 0 {
		clientID = "XXX"
	}
	var deviceInfo DeviceNode
	InitDeviceNode(&deviceInfo, nil, -1, clientID, "", "", "")
	err = json.Unmarshal(body, &deviceInfo)

	if err != nil {
		return nil, err
	}
	return &deviceInfo, nil
}

// ---------------------------------------------------------------------------

// Update reads the new value
func (node *DeviceNode) Update() error {
	httpClient := &http.Client{Timeout: 3 * time.Second}
	r, err := httpClient.Get(node.URLMeasure)

	if err != nil {
		return err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("Could not update device, status was '%d', message was '%s'!", r.StatusCode, body)
	}

	err = json.Unmarshal(body, &node.DeviceValue)

	if err != nil {
		return err
	}
	return nil
}

// ###########################################################################

// SetSequence ...
func (node *DeviceNode) SetSequence(s int) { node.sequence = s }

// ---------------------------------------------------------------------------

// GetSequence ...
func (node *DeviceNode) GetSequence() int { return node.sequence }

// ---------------------------------------------------------------------------

// IsNew ...
func (node *DeviceNode) IsNew() bool { return node.isNew }

// ---------------------------------------------------------------------------

// ClearNew ...
func (node *DeviceNode) ClearNew() {
	node.isNew = false
}

// ---------------------------------------------------------------------------

// GetURLDevice ...
func (node *DeviceNode) GetURLDevice() string { return node.URLDevice }

// ---------------------------------------------------------------------------

// GetURLMeasure ...
func (node *DeviceNode) GetURLMeasure() string { return node.URLMeasure }

// ---------------------------------------------------------------------------

// GetURLStatus ...
func (node *DeviceNode) GetURLStatus() string { return node.URLStatus }
