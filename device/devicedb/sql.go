package devicedb

// ###########################################################################
// ###########################################################################
// Database
// ###########################################################################
// ###########################################################################

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/com-gft-tsbo-source/go-common/device/implementation/devicemeasure"

	_ "github.com/boltdb/bolt"
	_ "github.com/lib/pq"
)

var tableStmtTmpl = "CREATE TABLE IF NOT EXISTS  %s (time VARCHAR, address VARCHAR, client VARCHAR, raw VARCHAR, unit VARCHAR, value VARCHAR)"
var updateStmtTmpl = "INSERT INTO %s (time, address, client, raw, unit, value) VALUES (<?1>, <?2>, <?3>, <?4>, <?5>, <?6>)"
var rePlaceholder = regexp.MustCompile(`(?:<\?(?P<idx>\d+)>)`)

// SQLConnection ...
type SQLConnection struct {
	Connection
	db         *sql.DB
	updateStmt *sql.Stmt
}

// ###########################################################################

// ---------------------------------------------------------------------------

// Open ...
func (con *SQLConnection) Open() error {

	con.mutex.Lock()
	defer con.mutex.Unlock()

	if con.db != nil {
		con.isOpen = true
		return nil
	}

	con.isOpen = false

	db, err := sql.Open(con.dbType, con.connectString)
	if err != nil {
		return err
	}

	_, err = db.Exec(con.tableStmtStr)
	if err != nil {
		db.Close()
		return err
	}

	con.updateStmt, err = db.Prepare(con.updateStmtStr)
	if err != nil {
		db.Close()
		return err
	}

	con.db = db
	con.isOpen = true
	return nil
}

// ---------------------------------------------------------------------------

// Close ...
func (con *SQLConnection) Close() error {

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

// ---------------------------------------------------------------------------

// AddDevice ...
func (con *SQLConnection) AddDevice(node devicemeasure.IDeviceMeasure) error {
	return nil
}

// ---------------------------------------------------------------------------

// Update ...
func (con *SQLConnection) Update(node devicemeasure.IDeviceMeasure) error {

	con.mutex.Lock()
	defer con.mutex.Unlock()

	var err error

	if !con.IsOpen() {
		return nil
	}

	_, err = con.updateStmt.Exec(node.GetStamp().Format("2006/01/02 15:04:05.000"), node.GetDeviceAddress(), con.Connection.clientName, fmt.Sprintf("%d", node.GetValue()), node.GetUnit(), node.GetFormatted())

	if err != nil {
		return err
	}

	return nil
}

// ---------------------------------------------------------------------------

// GetConnectString ...
func (con *SQLConnection) GetConnectString() string {
	return con.connectString
}

// ---------------------------------------------------------------------------

// IsOpen ...
func (con *SQLConnection) IsOpen() bool {
	return con.isOpen
}
