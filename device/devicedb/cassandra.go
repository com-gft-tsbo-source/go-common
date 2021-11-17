package devicedb

// ###########################################################################
// ###########################################################################
// Database
// ###########################################################################
// ###########################################################################

import (
	"fmt"
	"sync"

	"github.com/com-gft-tsbo-source/go-common/device/implementation/devicemeasure"

	"github.com/gocql/gocql"
)

// CassandraConnection ...
type CassandraConnection struct {
	Connection
	cluster *gocql.ClusterConfig
	session *gocql.Session
	path    string
}

var keyspaceStmtTmpl = "CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = { 'class': 'SimpleStrategy', 'replication_factor' : %d}"

// ###########################################################################

// NewDatabaseCassandra ...
func NewDatabaseCassandra(connectString string, tableName string, clientName string) (*CassandraConnection, bool) {
	path := connectString[6:]
	return &CassandraConnection{
		Connection{"Cassandra", path, tableName, clientName, "", "", false, &sync.Mutex{}},
		nil,
		nil,
		path,
	}, true
}

// Open ...
func (con *CassandraConnection) Open() error {
	var err error

	con.mutex.Lock()
	defer con.mutex.Unlock()

	if con.session != nil {
		con.isOpen = true
		return nil
	}

	con.cluster = nil
	con.session = nil
	con.isOpen = false
	fmt.Println(con.connectString)
	cluster := gocql.NewCluster(con.connectString)
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()

	if err != nil {
		con.cluster = nil
		con.session = nil
		return err
	}

	keyspaceStmt := fmt.Sprintf(keyspaceStmtTmpl, con.tableName, 1)
	err = session.Query(keyspaceStmt).Exec()

	if err != nil {
		return err
	}

	con.isOpen = true

	con.cluster = cluster
	con.session = session
	return nil
}

// ---------------------------------------------------------------------------

// Close ...
func (con *CassandraConnection) Close() error {

	con.mutex.Lock()
	defer con.mutex.Unlock()
	con.isOpen = false

	if con.session == nil {
		return nil
	}

	con.session.Close()
	con.cluster = nil
	con.session = nil
	return nil
}

// AddDevice ...
func (con *CassandraConnection) AddDevice(node devicemeasure.IDeviceMeasure) error {
	// con.mutex.Lock()
	// defer con.mutex.Unlock()

	// var err error

	// if !con.IsOpen() {
	// 	return nil
	// }

	// err = con.db.Update(func(tx *cassandra.Tx) error {
	// 	root, err := tx.CreateBucketIfNotExists([]byte(con.tableName))
	// 	if err != nil {
	// 		return fmt.Errorf("could not create root bucket: %v", err)
	// 	}
	// 	bucket, err := root.CreateBucketIfNotExists([]byte(con.clientName))
	// 	if err != nil {
	// 		return fmt.Errorf("could not create client bucket: %v", err)
	// 	}
	// 	bucket, err = bucket.CreateBucketIfNotExists([]byte(node.GetDeviceAddress()))
	// 	if err != nil {
	// 		return fmt.Errorf("could not create device bucket: %v", err)
	// 	}
	// 	return nil
	// })
	// return err
	return nil
}

// ---------------------------------------------------------------------------

// Update ...
func (con *CassandraConnection) Update(node devicemeasure.IDeviceMeasure) error {

	// con.mutex.Lock()
	// defer con.mutex.Unlock()

	// var err error

	// if !con.IsOpen() {
	// 	return nil
	// }

	// err = con.db.Update(func(tx *cassandra.Tx) error {
	// 	err := tx.Bucket([]byte(con.tableName)).Bucket([]byte(con.clientName)).Bucket([]byte(node.GetDeviceAddress())).Put([]byte(node.GetStamp().Format("2006/01/02 15:04:05.000")), []byte(node.GetFormatted()))
	// 	if err != nil {
	// 		return fmt.Errorf("could not insert value: %v", err)
	// 	}
	// 	return nil
	// })

	// if err != nil {
	// 	return err
	// }

	return nil
}

// ---------------------------------------------------------------------------

// GetConnectString ...
func (con *CassandraConnection) GetConnectString() string {
	return con.connectString
}

// ---------------------------------------------------------------------------

// IsOpen ...
func (con *CassandraConnection) IsOpen() bool {
	return con.isOpen
}
