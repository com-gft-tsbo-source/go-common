package devicedb

// ###########################################################################
// ###########################################################################
// Database
// ###########################################################################
// ###########################################################################

import (
	"fmt"
	"sync"
)

// NilConnection ...
type NilConnection struct {
	SQLConnection
}

// ###########################################################################

// NewDatabasePostgres ...
func NewDatabaseNil(connectString string, tableName string, clientName string) *NilConnection {
	updateStmtStr := rePlaceholder.ReplaceAllString(updateStmtTmpl, `$$$1`)
	updateStmtStr = fmt.Sprintf(updateStmtStr, tableName)
	tableStmtStr := fmt.Sprintf(tableStmtTmpl, tableName)

	return &NilConnection{
		SQLConnection{Connection{"nil", connectString, tableName, clientName, tableStmtStr, updateStmtStr, false, &sync.Mutex{}}, nil, nil},
	}
}
