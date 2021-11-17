package devicedb

// ###########################################################################
// ###########################################################################
// Database
// ###########################################################################
// ###########################################################################

import (
	"fmt"
	"sync"
	"time"

	"github.com/com-gft-tsbo-source/go-common/device/implementation/devicemeasure"

	"github.com/boltdb/bolt"
)

// BoltConnection ...
type BoltConnection struct {
	Connection
	db   *bolt.DB
	path string
}

// ###########################################################################

// NewDatabaseBolt ...
func NewDatabaseBolt(connectString string, tableName string, clientName string) (*BoltConnection, bool) {
	path := connectString[7:]
	return &BoltConnection{
		Connection{"Bolt", path, tableName, clientName, "", "", false, &sync.Mutex{}},
		nil,
		path,
	}, true
}

// Open ...
func (con *BoltConnection) Open() error {

	con.mutex.Lock()
	defer con.mutex.Unlock()

	if con.db != nil {
		con.isOpen = true
		return nil
	}

	con.isOpen = false

	db, err := bolt.Open(con.connectString, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	con.db = db
	con.isOpen = true

	return nil
}

// ---------------------------------------------------------------------------

// Close ...
func (con *BoltConnection) Close() error {

	con.mutex.Lock()
	defer con.mutex.Unlock()

	con.isOpen = false

	if con.db == nil {
		return nil
	}

	var err error

	err = con.db.Close()
	con.db = nil
	return err
}

// AddDevice ...
func (con *BoltConnection) AddDevice(node devicemeasure.IDeviceMeasure) error {

	con.mutex.Lock()
	defer con.mutex.Unlock()

	var err error

	if !con.IsOpen() {
		return nil
	}

	err = con.db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte(con.tableName))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		bucket, err := root.CreateBucketIfNotExists([]byte(con.clientName))
		if err != nil {
			return fmt.Errorf("could not create client bucket: %v", err)
		}
		bucket, err = bucket.CreateBucketIfNotExists([]byte(node.GetDeviceAddress()))
		if err != nil {
			return fmt.Errorf("could not create device bucket: %v", err)
		}
		return nil
	})
	return err
}

// ---------------------------------------------------------------------------

// Update ...
func (con *BoltConnection) Update(node devicemeasure.IDeviceMeasure) error {

	con.mutex.Lock()
	defer con.mutex.Unlock()

	var err error

	if !con.IsOpen() {
		return nil
	}

	err = con.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte(con.tableName)).Bucket([]byte(con.clientName)).Bucket([]byte(node.GetDeviceAddress())).Put([]byte(node.GetStamp().Format("2006/01/02 15:04:05.000")), []byte(node.GetFormatted()))
		if err != nil {
			return fmt.Errorf("could not insert value: %v", err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// ---------------------------------------------------------------------------

// GetConnectString ...
func (con *BoltConnection) GetConnectString() string {
	return con.connectString
}

// ---------------------------------------------------------------------------

// IsOpen ...
func (con *BoltConnection) IsOpen() bool {
	return con.isOpen
}
