package devicedb

// ###########################################################################
// ###########################################################################
// Database
// ###########################################################################
// ###########################################################################

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
)

// PostgresqlConnection ...
type PostgresqlConnection struct {
	SQLConnection
	user     string
	password string
	host     string
	port     int
	database string
}

// ###########################################################################

// PGConnectRegexp ...
var PGConnectRegexp = regexp.MustCompile(`pg://(?P<user>[^:]+):(?P<pass>[^@]+)@(?P<host>[^:]+):(?P<port>\d+)/(?P<db>.+)`)

// NewDatabasePostgres ...
func NewDatabasePostgres(connectString string, tableName string, clientName string) (*PostgresqlConnection, bool) {
	m := PGConnectRegexp.FindStringSubmatch(connectString)

	if m == nil {
		return nil, false
	}

	port, err := strconv.Atoi(m[4])
	if err != nil {
		port = 5432
	}

	updateStmtStr := rePlaceholder.ReplaceAllString(updateStmtTmpl, `$$$1`)
	updateStmtStr = fmt.Sprintf(updateStmtStr, tableName)
	tableStmtStr := fmt.Sprintf(tableStmtTmpl, tableName)
	connectString = fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable connect_timeout=3", m[1], m[2], m[3], port, m[5])

	return &PostgresqlConnection{
		SQLConnection{Connection{"postgres", connectString, tableName, clientName, tableStmtStr, updateStmtStr, false, &sync.Mutex{}}, nil, nil},
		m[1],
		m[2],
		m[3],
		port,
		m[5],
	}, true
}
