package devicedb

// ###########################################################################
// ###########################################################################
// Database
// ###########################################################################
// ###########################################################################

import (
	"strings"
	"sync"

	"github.com/com-gft-tsbo-source/go-common/device/implementation/devicemeasure"

	_ "github.com/boltdb/bolt"
	_ "github.com/lib/pq"
)

// Connection ...
type Connection struct {
	dbType        string
	connectString string
	tableName     string
	clientName    string
	tableStmtStr  string
	updateStmtStr string
	isOpen        bool
	mutex         *sync.Mutex
}

// IConnection ...
type IConnection interface {
	Open() error
	Close() error
	AddDevice(devicemeasure.IDeviceMeasure) error
	Update(devicemeasure.IDeviceMeasure) error
	GetConnectString() string
	IsOpen() bool
}

// ###########################################################################

// NewDatabase ...
func NewDatabase(connectString string, tableName string, clientName string) IConnection {

	var retVal IConnection
	var isOk bool
	if strings.HasPrefix(connectString, "sql://") {
		retVal, isOk = NewDatabaseSqlite(connectString, tableName, clientName)
	} else if strings.HasPrefix(connectString, "pg://") {
		retVal, isOk = NewDatabasePostgres(connectString, tableName, clientName)
	} else if strings.HasPrefix(connectString, "bolt://") {
		retVal, isOk = NewDatabaseBolt(connectString, tableName, clientName)
	} else if strings.HasPrefix(connectString, "cql://") {
		retVal, isOk = NewDatabaseCassandra(connectString, tableName, clientName)
	}

	if isOk {
		return retVal
	}

	return NewDatabaseNil(connectString, tableName, clientName)
}
