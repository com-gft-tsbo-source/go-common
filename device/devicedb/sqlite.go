// +build sqlite

package devicedb

// ###########################################################################
// ###########################################################################
// Database
// ###########################################################################
// ###########################################################################

import (
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type ConnectionSqlite struct {
	SQLConnection
	path string
}

// ###########################################################################

func NewDatabaseSqlite(connectString string, tableName string, clientName string) (*ConnectionSqlite, bool) {

	tableStmtStr := fmt.Sprintf(tableStmtTmpl, tableName)
	updateStmtStr := rePlaceholder.ReplaceAllString(updateStmtTmpl, "?")
	updateStmtStr = fmt.Sprintf(updateStmtStr, tableName)
	path := connectString[5:]
	return &ConnectionSqlite{
		SQLConnection{Connection{"sqlite3", path, tableName, clientName, tableStmtStr, updateStmtStr, false, &sync.Mutex{}}, nil, nil},
		path,
	}, true
}
